package application

import (
	"fmt"
	"rynx/internal/jira"
	"rynx/shared/configutils"
)

func (a *Application_struct) Application_ping_method() error {
	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}

	client := jira.Jira_client_const(cfg)
	res, err := client.Get_myself_method()
	if err != nil {
		return err
	}

	fmt.Printf("✅ Connected as %s (%s)\n", res.DisplayName, res.EmailAddress)
	return nil
}
