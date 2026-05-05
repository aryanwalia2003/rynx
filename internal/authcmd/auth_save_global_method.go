package authcmd

import (
	"rynx/internal/config"
	"rynx/shared/configutils"
	"rynx/shared/errors"
)

func (a *Auth_struct) save_global_auth_method(email, token string) error {
	path, err := configutils.Config_global_path_util()
	if err != nil {
		return err
	}

	cfg, _ := configutils.Config_load_util(path)
	if cfg == nil {
		cfg = &config.Config_struct{}
	}

	cfg.Jira.UserEmail = email
	cfg.Jira.Token = token
	cfg.Jira.AuthType = "api_token"

	if err := configutils.Config_save_util(path, cfg); err != nil {
		return errors.Error_wrap_util("failed to save global config", err)
	}
	return nil
}
