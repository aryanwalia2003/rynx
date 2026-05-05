package jira

import (
	"rynx/internal/config"
)

type Jira_client_struct struct {
	config *config.Config_struct
}

func (j *Jira_client_struct) Get_ticket_method(id string) (string, error) {
	return "Ticket info for " + id, nil
}
