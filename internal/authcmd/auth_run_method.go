package authcmd

import (
	"rynx/shared/promptutils"
)

func (a *Auth_struct) Auth_run_method() error {
	domain, _ := promptutils.Prompt_string_util("Jira Domain", "")
	email, _ := promptutils.Prompt_string_util("User Email", "")
	token, _ := promptutils.Prompt_string_util("API Token", "")

	if err := a.save_local_auth_method(domain, email); err != nil {
		return err
	}
	return a.save_global_auth_method(email, token)
}
