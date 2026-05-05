package detectutils

import (
	"rynx/internal/config"
	"strings"
)

func Jira_detect_util(gitCfg *config.Git_config_struct) *config.Jira_config_struct {
	parts := strings.Split(gitCfg.OriginURL, "/")
	if len(parts) > 0 {
		repoName := parts[len(parts)-1]
		repoName = strings.TrimSuffix(repoName, ".git")
		projKey := strings.ToUpper(repoName)
		return &config.Jira_config_struct{
			Domain:            "company.atlassian.net",
			DefaultProjectKey: projKey,
		}
	}

	return &config.Jira_config_struct{}
}
