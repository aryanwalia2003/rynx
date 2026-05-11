package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

// Jira_user_struct holds the currently authenticated user's info.
type Jira_user_struct struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

func (j *Jira_client_struct) Get_myself_method() (*Jira_user_struct, error) {
	domain := j.config.Jira.Domain
	url := fmt.Sprintf("https://%s/rest/api/3/myself", domain)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Error_wrap_util("ping request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ping failed with status: %d", resp.StatusCode)
	}

	var user Jira_user_struct
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errors.Error_wrap_util("failed to decode myself response", err)
	}
	return &user, nil
}
