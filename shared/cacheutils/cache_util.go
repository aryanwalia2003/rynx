package cacheutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"rynx/shared/configutils"
	"time"
)

type TicketDraftStruct struct {
	ProjectKey   string `json:"project_key"`
	IssueTypeID  string `json:"issue_type_id"`
	IssueType    string `json:"issue_type"`
	Summary      string `json:"summary"`
	Description  string `json:"description"`
	Attachment   string `json:"attachment"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Cache struct {
	BranchTickets  map[string]string `json:"branch_tickets"` // Branch name -> Ticket ID
	RecentProjects []string          `json:"recent_projects"`
	TicketDraft    *TicketDraftStruct `json:"ticket_draft,omitempty"`
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
		return &Cache{
			BranchTickets:  make(map[string]string),
			RecentProjects: make([]string, 0),
		}, nil
	}

	if cache.BranchTickets == nil {
		cache.BranchTickets = make(map[string]string)
	}
	if cache.RecentProjects == nil {
		cache.RecentProjects = make([]string, 0)
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
