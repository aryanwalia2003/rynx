package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

type jira_agile_issue_response_struct struct {
	Issues []struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Status  struct {
				Name string `json:"name"`
			} `json:"status"`
		} `json:"fields"`
	} `json:"issues"`
}

func (j *Jira_client_struct) Jira_get_active_sprint_tickets_method(projectKey string) ([]string, error) {
	// First, find the board for this project
	boardURL := fmt.Sprintf("https://%s/rest/agile/1.0/board?projectKeyOrId=%s", j.config.Jira.Domain, projectKey)
	req, _ := http.NewRequest("GET", boardURL, nil)
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("board fetch failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("board fetch api error: %d", resp.StatusCode)
	}

	var boardData struct {
		Values []struct {
			ID int `json:"id"`
		} `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&boardData); err != nil {
		return nil, errors.Error_wrap_util("failed to decode board response", err)
	}

	if len(boardData.Values) == 0 {
		return nil, fmt.Errorf("no board found for project %s", projectKey)
	}

	boardID := boardData.Values[0].ID

	// Now get issues from that board
	// The user wants to see their tickets on the active sprint board (matching the UI filter).
	issuesURL := fmt.Sprintf("https://%s/rest/agile/1.0/board/%d/issue?jql=sprint+in+openSprints()+AND+assignee=currentUser()", j.config.Jira.Domain, boardID)
	req, _ = http.NewRequest("GET", issuesURL, nil)
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("agile issues request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback: If openSprints() fails, just get board issues for the current user
		issuesURL = fmt.Sprintf("https://%s/rest/agile/1.0/board/%d/issue?jql=assignee=currentUser()", j.config.Jira.Domain, boardID)
		req, _ = http.NewRequest("GET", issuesURL, nil)
		req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("agile api error: %d", resp.StatusCode)
		}
	}

	var data jira_agile_issue_response_struct
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Error_wrap_util("failed to decode agile response", err)
	}

	results := []string{}
	for _, issue := range data.Issues {
		results = append(results, fmt.Sprintf("%s | %s [%s]", issue.Key, issue.Fields.Summary, issue.Fields.Status.Name))
	}
	return results, nil
}
