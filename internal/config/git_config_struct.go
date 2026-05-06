package config

type Git_config_struct struct {
	OriginURL    string `json:"origin_url"`
	UpstreamURL  string `json:"upstream_url,omitempty"`
	DefaultBranch string `json:"default_branch"`
	GitHubToken  string `json:"github_token,omitempty"`
}
