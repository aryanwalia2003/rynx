package initcmd

import (
	"rynx/internal/config"
	"rynx/shared/configutils"
	"rynx/shared/detectutils"
	"rynx/shared/errors"
	"rynx/shared/promptutils"
	"strings"
)

func (cmd *Init_struct) Init_run_method() error {
	cmd.logger.Logger_info_method("Running dev init...")

	gitCfg, err := detectutils.Git_detect_util()
	if err != nil {
		return errors.Error_wrap_util("git detection failed", err)
	}
	cmd.logger.Logger_info_method("Detected Git remote: " + gitCfg.OriginURL)

	jiraCfg := detectutils.Jira_detect_util(gitCfg)

	msg := "Enter Jira project key (or 'skip' to disable Jira integration)"
	userProj, err := promptutils.Prompt_string_util(msg, jiraCfg.DefaultProjectKey)
	if err != nil {
		return errors.Error_wrap_util("prompt failed", err)
	}

	if strings.ToLower(userProj) == "skip" {
		jiraCfg.DefaultProjectKey = ""
		jiraCfg.Domain = ""
		cmd.logger.Logger_info_method("Jira integration skipped.")
	} else {
		jiraCfg.DefaultProjectKey = userProj
		cmd.logger.Logger_info_method("Configured Jira project: " + jiraCfg.DefaultProjectKey)
	}

	cfg := &config.Config_struct{
		Git:  *gitCfg,
		Jira: *jiraCfg,
	}

	if err := configutils.Config_save_util(".devrc", cfg); err != nil {
		return errors.Error_wrap_util("failed to save config", err)
	}

	cmd.logger.Logger_info_method("Configuration saved to .devrc")
	return nil
}
