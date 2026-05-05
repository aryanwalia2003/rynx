package configutils

import (
	"encoding/json"
	"os"
	"rynx/internal/config"
	"rynx/shared/errors"
)

func Config_load_util(path string) (*config.Config_struct, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &config.Config_struct{}, nil
		}
		return nil, errors.Error_wrap_util("failed to read config file", err)
	}

	var c config.Config_struct
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, errors.Error_wrap_util("failed to parse config JSON", err)
	}

	return &c, nil
}
