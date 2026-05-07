package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

// Jira_transition_struct represents a single available transition on a ticket.
type Jira_transition_struct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Fields holds the schema for any required fields on this transition's screen.
	Fields map[string]Jira_transition_field_struct `json:"fields"`
}

// Jira_transition_field_struct describes a single required field on a transition screen.
type Jira_transition_field_struct struct {
	Required bool   `json:"required"`
	Name     string `json:"name"`
	Schema   struct {
		Type   string `json:"type"`
		System string `json:"system"`
		Items  string `json:"items"`
	} `json:"schema"`
	AllowedValues []Jira_allowed_value_struct `json:"allowedValues"`
}

// Jira_allowed_value_struct represents an option in a dropdown/select field.
type Jira_allowed_value_struct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Jira_get_transitions_method fetches all available transitions for a given issue key.
func (j *Jira_client_struct) Jira_get_transitions_method(key string) ([]Jira_transition_struct, error) {
	path := fmt.Sprintf("https://%s/rest/api/3/issue/%s/transitions?expand=transitions.fields", j.config.Jira.Domain, key)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, errors.Error_wrap_util("failed to build transitions request", err)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("transitions fetch failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transitions error: HTTP %d for %s", resp.StatusCode, key)
	}

	var result struct {
		Transitions []Jira_transition_struct `json:"transitions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Error_wrap_util("failed to decode transitions", err)
	}
	return result.Transitions, nil
}
