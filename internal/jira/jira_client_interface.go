package jira

import (
	"rynx/internal/config"
)

type Jira_client_interface interface {
	Get_ticket_method(id string) (string, error)
	Jira_search_method(jql string) ([]string, error)
	Get_myself_method() (*Jira_user_struct, error)
	Jira_get_issue_detail_method(key string) (*Jira_issue_result_struct, error)
	Jira_get_active_sprint_tickets_method(projectKey string) ([]string, error)
	Jira_get_my_tickets_method(projectKey string) ([]string, error)
	Jira_get_transitions_method(key string) ([]Jira_transition_struct, error)
	Jira_do_transition_method(key string, transitionID string, fields Jira_transition_fields_map) error
	Jira_get_projects_method() ([]Jira_project_struct, error)
	Jira_get_issue_types_method(projectKey string) ([]Jira_issue_type_struct, error)
	Jira_create_issue_method(fields map[string]interface{}) (string, error)
	Jira_add_attachment_method(issueKey string, filePath string) error
}

func Jira_client_const(cfg *config.Config_struct) Jira_client_interface {
	return &Jira_client_struct{
		config: cfg,
	}
}
