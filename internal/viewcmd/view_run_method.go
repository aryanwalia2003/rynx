package viewcmd

import (
	"fmt"
	"rynx/internal/jira"
	"rynx/shared/configutils"
	"strings"
)

const (
	ansi_bold       = "\033[1m"
	ansi_reset      = "\033[0m"
	ansi_cyan       = "\033[36m"
	ansi_green      = "\033[32m"
	ansi_yellow     = "\033[33m"
)

func (v *View_struct) View_run_method(args []string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a ticket ID (e.g. dev view ZZDPA-123)")
		return nil
	}

	ticketID := strings.ToUpper(args[0])

	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}

	client := jira.Jira_client_const(cfg)

	fmt.Printf("Fetching details for %s...\n", ticketID)

	detail, err := client.Jira_get_issue_detail_method(ticketID)
	if err != nil {
		return err
	}

	// Print formatted output
	fmt.Print("\033[2J\033[H") // Clear screen
	fmt.Printf("%s%s[%s]%s %s\n", ansi_bold, ansi_cyan, detail.Key, ansi_reset, detail.Summary)
	fmt.Printf("%sStatus:%s %s%s%s\n", ansi_bold, ansi_reset, ansi_green, detail.Status, ansi_reset)
	
	if detail.Assignee != "" {
		fmt.Printf("%sAssignee:%s %s\n", ansi_bold, ansi_reset, detail.Assignee)
	}
	if detail.Reporter != "" {
		fmt.Printf("%sReporter:%s %s\n", ansi_bold, ansi_reset, detail.Reporter)
	}
	
	fmt.Printf("\n%sDescription:%s\n", ansi_bold, ansi_reset)
	fmt.Println(strings.Repeat("-", 40))
	
	desc := strings.TrimSpace(detail.Description)
	if desc == "" {
		fmt.Println("No description provided.")
	} else {
		fmt.Println(desc)
	}
	fmt.Println(strings.Repeat("-", 40))

	return nil
}
