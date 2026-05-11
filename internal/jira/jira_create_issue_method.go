package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"rynx/shared/errors"
)

// Jira_create_issue_payload is the body sent to POST /rest/api/3/issue.
// Description is encoded in Atlassian Document Format (ADF).
type Jira_create_issue_payload struct {
	Fields map[string]interface{} `json:"fields"`
}

// Jira_create_issue_method creates a new Jira issue and returns its key (e.g. "PROJ-42").
// Caller is responsible for building the fields map correctly.
func (j *Jira_client_struct) Jira_create_issue_method(fields map[string]interface{}) (string, error) {
	payload := Jira_create_issue_payload{Fields: fields}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Error_wrap_util("failed to marshal issue payload", err)
	}

	url := fmt.Sprintf("https://%s/rest/api/3/issue", j.config.Jira.Domain)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", errors.Error_wrap_util("failed to build create issue request", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Error_wrap_util("create issue request failed", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create issue API error: HTTP %d — %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", errors.Error_wrap_util("failed to decode create issue response", err)
	}
	return result.Key, nil
}

// Jira_add_attachment_method uploads a file to an existing issue.
// The Jira API requires multipart/form-data with an X-Atlassian-Token: no-check header.
func (j *Jira_client_struct) Jira_add_attachment_method(issueKey string, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return errors.Error_wrap_util("failed to open attachment file", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return errors.Error_wrap_util("failed to create form file", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return errors.Error_wrap_util("failed to copy file to form", err)
	}
	writer.Close()

	url := fmt.Sprintf("https://%s/rest/api/3/issue/%s/attachments", j.config.Jira.Domain, issueKey)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return errors.Error_wrap_util("failed to build attachment request", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Atlassian-Token", "no-check") // Required by Jira to bypass XSRF check
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Error_wrap_util("attachment upload request failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("attachment upload error: HTTP %d — %s", resp.StatusCode, string(body))
	}
	return nil
}

// Build_adf_description_util converts a plain-text multi-line string into
// the minimal Atlassian Document Format (ADF) required by the Jira v3 API.
func Build_adf_description_util(text string) map[string]interface{} {
	return map[string]interface{}{
		"type":    "doc",
		"version": 1,
		"content": []map[string]interface{}{
			{
				"type": "paragraph",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": text,
					},
				},
			},
		},
	}
}
