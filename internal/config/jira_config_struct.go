package config

type Jira_config_struct struct {
	Domain            string   `json:"domain"`
	DefaultProjectKey string   `json:"default_project_key"`
	AllowedProjectKeys []string `json:"allowed_project_keys,omitempty"`
	AuthType          string   `json:"auth_type,omitempty"`
	UserEmail         string   `json:"user_email,omitempty"`
	Token             string   `json:"token,omitempty"`
}
