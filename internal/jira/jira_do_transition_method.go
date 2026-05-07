package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rynx/shared/errors"
)

// Jira_transition_fields_map holds field values to submit with a transition.
type Jira_transition_fields_map map[string]interface{}

// Jira_do_transition_method executes a transition on a Jira issue.
// fields can be nil if the transition requires no extra screen fields.
func (j *Jira_client_struct) Jira_do_transition_method(key string, transitionID string, fields Jira_transition_fields_map) error {
	payload := map[string]interface{}{
		"transition": map[string]string{"id": transitionID},
	}
	if len(fields) > 0 {
		payload["fields"] = fields
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return errors.Error_wrap_util("failed to marshal transition payload", err)
	}

	path := fmt.Sprintf("https://%s/rest/api/3/issue/%s/transitions", j.config.Jira.Domain, key)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(body))
	if err != nil {
		return errors.Error_wrap_util("failed to build transition POST request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(j.config.Jira.UserEmail, j.config.Jira.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Error_wrap_util("transition POST failed", err)
	}
	defer resp.Body.Close()

	// Jira returns 204 No Content on success
	if resp.StatusCode != http.StatusNoContent {
		var errBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errBody)
		msg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if msgs, ok := errBody["errorMessages"].([]interface{}); ok && len(msgs) > 0 {
			msg = fmt.Sprintf("%v", msgs[0])
		}
		return fmt.Errorf("transition failed: %s", msg)
	}
	return nil
}
