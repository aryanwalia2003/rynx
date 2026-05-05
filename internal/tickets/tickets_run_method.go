package tickets

import (
	"fmt"
	"rynx/internal/jira"
	"rynx/shared/configutils"
	"rynx/shared/flagutils"
)

func (t *Tickets_struct) Tickets_run_method(args []string) error {
	flags := flagutils.Flag_parser_util(args)
	jql := jira.Jql_builder_util(flags)
	
	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}

	client := jira.Jira_client_const(cfg)
	results, err := client.Jira_search_method(jql)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No tickets found matching your filters.")
		return nil
	}

	for _, res := range results {
		fmt.Println(res)
	}
	return nil
}
