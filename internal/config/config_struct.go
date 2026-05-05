package config

type Config_struct struct {
	Git  Git_config_struct  `json:"git"`
	Jira Jira_config_struct `json:"jira"`
}
