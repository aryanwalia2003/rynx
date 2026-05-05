package configutils

import (
	"encoding/json"
	"os"
	"rynx/internal/config"
	"rynx/shared/errors"
)

func Config_save_util(path string, c *config.Config_struct) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.Error_wrap_util("failed to serialize config", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return errors.Error_wrap_util("failed to write config file", err)
	}

	return nil
}
