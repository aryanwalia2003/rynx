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
	results, err := client.Jira_get_active_sprint_tickets_method(projectKey)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No tickets found in the active sprint.")
		return nil
	}

	// If --status flag provided, do a simple substring filter
	if statusFilter := strings.ToLower(flags["status"]); statusFilter != "" {
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

	// No flag: launch interactive fuzzy picker
	selected, ok := ttyutils.Tty_picker_util(results, "Filter tickets")
	fmt.Print("\033[2J\033[H")
	if ok {
		fmt.Println(selected)
	}
	return nil
}
