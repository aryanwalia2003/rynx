package tickets

import (
	"fmt"
	"rynx/internal/jira"
	"rynx/shared/configutils"
	"rynx/shared/flagutils"
	"rynx/shared/ttyutils"
	"strings"
)

func (t *Tickets_struct) Tickets_run_method(args []string) error {
	flags := flagutils.Flag_parser_util(args)
	projectKey := flags["project"]
	if projectKey == "" {
		fmt.Println("Please provide a project key using --project")
		return nil
	}

	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}

	client := jira.Jira_client_const(cfg)
	results, err := client.Jira_get_my_tickets_method(projectKey)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No tickets found.")
		return nil
	}

	_, hasStatusFlag := flags["status"]
	statusVal := flags["status"]

	if hasStatusFlag {
		if statusVal == "" {
			// Extract all unique statuses
			statusSet := make(map[string]bool)
			for _, r := range results {
				parts := strings.Split(r, "[")
				if len(parts) > 1 {
					status := strings.TrimRight(parts[len(parts)-1], "]")
					statusSet[status] = true
				}
			}

			var statuses []string
			for s := range statusSet {
				statuses = append(statuses, s)
			}

			if len(statuses) == 0 {
				fmt.Println("No statuses found.")
				return nil
			}

			selectedStatus, ok := ttyutils.Tty_picker_util(statuses, "Select Status")
			if !ok {
				return nil
			}

			fmt.Print("\033[2J\033[H") // Clear screen
			fmt.Printf("Tickets in status: %s\n\n", selectedStatus)
			for _, r := range results {
				parts := strings.Split(r, "[")
				if len(parts) > 1 {
					status := strings.TrimRight(parts[len(parts)-1], "]")
					if status == selectedStatus {
						fmt.Println(r)
					}
				}
			}
			return nil
		} else {
			// If --status provided a value, do a simple substring filter
			statusFilter := strings.ToLower(statusVal)
			for _, r := range results {
				parts := strings.Split(r, "[")
				if len(parts) > 1 {
					status := strings.ToLower(strings.TrimRight(parts[len(parts)-1], "]"))
					if strings.Contains(status, statusFilter) {
						fmt.Println(r)
					}
				}
			}
			return nil
		}
	}

	// No flag: print all my tickets
	for _, r := range results {
		fmt.Println(r)
	}
	return nil
}
