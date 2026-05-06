package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

func (j *Jira_client_struct) Jira_get_my_tickets_method(projectKey string) ([]string, error) {
	jql := fmt.Sprintf("project=%s AND assignee=currentUser() ORDER BY updated DESC", projectKey)

	body, _ := json.Marshal(map[string]interface{}{
		"jql":        jql,
		"maxResults": 50,
		"fields":     []string{"summary", "status"},
	})

	path := fmt.Sprintf("https://%s/rest/api/3/search/jql", j.config.Jira.Domain)
	req, _ := http.NewRequest("POST", path, bytes.NewReader(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("my tickets request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jira api error: %d", resp.StatusCode)
	}

	var data struct {
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

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Error_wrap_util("failed to decode search response", err)
	}

	results := []string{}
	for _, issue := range data.Issues {
		results = append(results, fmt.Sprintf("%s | %s [%s]", issue.Key, issue.Fields.Summary, issue.Fields.Status.Name))
	}
	return results, nil
}
