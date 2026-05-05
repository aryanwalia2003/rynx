package configutils

import (
	"os"
	"path/filepath"
	"rynx/shared/errors"
)

func Config_global_path_util() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Error_wrap_util("failed to get home dir", err)
	}
	return filepath.Join(home, ".devrc"), nil
}
