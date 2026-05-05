package jira

import (
	"strings"
)

func (j *Jira_client_struct) extract_query_from_jql_method(jql string) string {
	parts := strings.Split(jql, "AND")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "project") {
			parts := strings.Split(part, "=")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}
