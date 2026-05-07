package jira

import (
	"fmt"
	"strings"
)

func Jql_builder_util(filters map[string]string) string {
	var parts []string

	if proj, ok := filters["project"]; ok {
		parts = append(parts, fmt.Sprintf("project = %s", proj))
	} else {
		parts = append(parts, "assignee = currentUser()")
	}

	if status, ok := filters["status"]; ok {
		parts = append(parts, fmt.Sprintf("status = '%s'", status))
	}

	if len(parts) == 0 {
		return "assignee = currentUser()"
	}

	return strings.Join(parts, " AND ")
}
