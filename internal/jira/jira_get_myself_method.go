package jira

import (
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

func (j *Jira_client_struct) Get_myself_method() (string, error) {
	domain := j.config.Jira.Domain
	url := fmt.Sprintf("https://%s/rest/api/3/myself", domain)

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Error_wrap_util("ping request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ping failed with status: %d", resp.StatusCode)
	}

	return "Connection successful!", nil
}
