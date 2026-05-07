package jira

import (
	"rynx/internal/config"
)

type Jira_client_interface interface {
	Get_ticket_method(id string) (string, error)
	Jira_search_method(jql string) ([]string, error)
	Get_myself_method() (string, error)
	Jira_get_issue_detail_method(key string) (*Jira_issue_result_struct, error)
	Jira_get_active_sprint_tickets_method(projectKey string) ([]string, error)
	Jira_get_my_tickets_method(projectKey string) ([]string, error)
	Jira_get_transitions_method(key string) ([]Jira_transition_struct, error)
	Jira_do_transition_method(key string, transitionID string, fields Jira_transition_fields_map) error
}

func Jira_client_const(cfg *config.Config_struct) Jira_client_interface {
	return &Jira_client_struct{
		config: cfg,
	}
}
