package jira

import (
	"rynx/internal/config"
)

type Jira_client_interface interface {
	Get_ticket_method(id string) (string, error)
	Jira_search_method(jql string) ([]string, error)
	Get_myself_method() (string, error)
}

func Jira_client_const(cfg *config.Config_struct) Jira_client_interface {
	return &Jira_client_struct{
		config: cfg,
	}
}
