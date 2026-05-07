package configutils

import (
	"rynx/internal/config"
)

func Config_load_merged_util() (*config.Config_struct, error) {
	local, err := Config_load_util(".devrc")
	if err != nil {
		return nil, err
	}

	path, _ := Config_global_path_util()
	global, _ := Config_load_util(path)

	if global != nil {
		local.Jira.Token = global.Jira.Token
		if local.Jira.UserEmail == "" {
			local.Jira.UserEmail = global.Jira.UserEmail
		}
	}

	return local, nil
}
