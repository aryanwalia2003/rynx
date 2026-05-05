package detectutils

import (
	"os/exec"
	"rynx/internal/config"
	"rynx/shared/errors"
	"strings"
)

func Git_detect_util() (*config.Git_config_struct, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Error_wrap_util("failed to get git remote", err)
	}

	origin := strings.TrimSpace(string(out))
	if origin == "" {
		return nil, errors.Error_wrap_util("empty origin URL", nil)
	}

	return &config.Git_config_struct{
		OriginURL:    origin,
		DefaultBranch: "main",
	}, nil
}
