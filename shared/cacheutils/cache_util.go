package cacheutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"rynx/shared/configutils"
)

type Cache struct {
	BranchTickets map[string]string `json:"branch_tickets"` // Branch name -> Ticket ID
}

func getCachePath() (string, error) {
	globalPath, err := configutils.Config_global_path_util()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(globalPath)
	return filepath.Join(dir, "cache.json"), nil
}

func Load_cache_util() (*Cache, error) {
	path, err := getCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Cache{BranchTickets: make(map[string]string)}, nil
		}
		return nil, err
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return &Cache{BranchTickets: make(map[string]string)}, nil
	}

	if cache.BranchTickets == nil {
		cache.BranchTickets = make(map[string]string)
	}

	return &cache, nil
}

func Save_cache_util(cache *Cache) error {
	path, err := getCachePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
