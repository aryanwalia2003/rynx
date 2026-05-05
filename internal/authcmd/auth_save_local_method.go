package authcmd

import (
	"rynx/shared/configutils"
	"rynx/shared/errors"
)

func (a *Auth_struct) save_local_auth_method(domain, email string) error {
	cfg, err := configutils.Config_load_util(".devrc")
	if err != nil {
		return errors.Error_wrap_util("failed to load local config", err)
	}

	cfg.Jira.Domain = domain
	cfg.Jira.UserEmail = email
	cfg.Jira.AuthType = "api_token"

	if err := configutils.Config_save_util(".devrc", cfg); err != nil {
		return errors.Error_wrap_util("failed to save local config", err)
	}
	return nil
}
