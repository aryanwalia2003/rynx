package tickets

import (
	"fmt"
	"rynx/internal/jira"
	"rynx/shared/configutils"
	"rynx/shared/flagutils"
	"strings"
)

func (t *Tickets_struct) Tickets_run_method(args []string) error {
	flags := flagutils.Flag_parser_util(args)
	projectKey := flags["project"]
	if projectKey == "" {
		fmt.Println("Please provide a project key using --project")
		return nil
	}

	statusFilter := strings.ToLower(flags["status"])

	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}

	client := jira.Jira_client_const(cfg)
	
	// Fetch tickets using Agile API to match the active sprint board
	results, err := client.Jira_get_active_sprint_tickets_method(projectKey)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No tickets found in the active sprint.")
		return nil
	}

	printed := 0
	for _, result := range results {
		// results are in format "KEY | Summary [Status]"
		// We extract the status to filter it
		if statusFilter != "" {
			parts := strings.Split(result, "[")
			if len(parts) > 1 {
				statusPart := strings.TrimRight(parts[len(parts)-1], "]")
				if !strings.Contains(strings.ToLower(statusPart), statusFilter) {
					continue
				}
			}
		}
		fmt.Println(result)
		printed++
	}

	if printed == 0 {
		fmt.Println("No tickets matched your filters.")
	}
	return nil
}
