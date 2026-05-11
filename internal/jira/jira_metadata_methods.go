package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

// Jira_project_struct represents a Jira project.
type Jira_project_struct struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Jira_issue_type_struct represents an issue type within a project.
type Jira_issue_type_struct struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Subtask     bool   `json:"subtask"`
}

// Jira_get_projects_method returns all projects the authenticated user can access.
func (j *Jira_client_struct) Jira_get_projects_method() ([]Jira_project_struct, error) {
	url := fmt.Sprintf("https://%s/rest/api/3/project/search?maxResults=100&orderBy=lastIssueUpdatedTime", j.config.Jira.Domain)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Error_wrap_util("failed to build projects request", err)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("projects fetch failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("projects API error: HTTP %d", resp.StatusCode)
	}

	var result struct {
		Values []Jira_project_struct `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Error_wrap_util("failed to decode projects response", err)
	}
	return result.Values, nil
}

// Jira_get_issue_types_method returns all non-subtask issue types for a given project key.
func (j *Jira_client_struct) Jira_get_issue_types_method(projectKey string) ([]Jira_issue_type_struct, error) {
	url := fmt.Sprintf("https://%s/rest/api/3/project/%s", j.config.Jira.Domain, projectKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Error_wrap_util("failed to build issue types request", err)
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("issue types fetch failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("issue types API error: HTTP %d for project %s", resp.StatusCode, projectKey)
	}

	var result struct {
		IssueTypes []Jira_issue_type_struct `json:"issueTypes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Error_wrap_util("failed to decode issue types response", err)
	}

	// Filter out subtasks — keep things simple
	var filtered []Jira_issue_type_struct
	for _, it := range result.IssueTypes {
		if !it.Subtask {
			filtered = append(filtered, it)
		}
	}
	return filtered, nil
}
