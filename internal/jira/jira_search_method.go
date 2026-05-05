package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"rynx/shared/errors"
)

type jira_picker_response_struct struct {
	Sections []struct {
		Issues []struct {
			Key        string `json:"key"`
			SummaryText string `json:"summaryText"`
		} `json:"issues"`
	} `json:"sections"`
}

func (j *Jira_client_struct) Jira_search_method(jql string) ([]string, error) {
	query := j.extract_query_from_jql_method(jql)
	domain := j.config.Jira.Domain
	path := fmt.Sprintf(
		"https://%s/rest/api/3/issue/picker?currentJQL=%s&query=%s",
		domain,
		url.QueryEscape(jql),
		url.QueryEscape(query),
	)

	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("picker request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jira api error: %d", resp.StatusCode)
	}

	var data jira_picker_response_struct
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, errors.Error_wrap_util("failed to decode picker response", err)
	}

	seen := map[string]bool{}
	results := []string{}
	for _, section := range data.Sections {
		for _, issue := range section.Issues {
			if !seen[issue.Key] {
				seen[issue.Key] = true
				results = append(results, fmt.Sprintf("%s | %s", issue.Key, issue.SummaryText))
			}
		}
	}
	return results, nil
}
