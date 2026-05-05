package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

type jira_issue_detail_struct struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

func (j *Jira_client_struct) fetch_issue_detail_method(key string) (*jira_issue_detail_struct, error) {
	path := fmt.Sprintf("https://%s/rest/api/3/issue/%s?fields=summary,status", j.config.Jira.Domain, key)
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("issue detail fetch failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("issue detail error: %d for %s", resp.StatusCode, key)
	}

	var detail jira_issue_detail_struct
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, errors.Error_wrap_util("failed to decode issue detail", err)
	}
	return &detail, nil
}
