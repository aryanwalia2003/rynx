package createcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"rynx/internal/jira"
	"rynx/internal/startcmd"
	"rynx/shared/cacheutils"
	"rynx/shared/configutils"
	"rynx/shared/errors"
	"rynx/shared/flagutils"
	"rynx/shared/promptutils"
	"rynx/shared/ttyutils"
)

const (
	ansi_bold   = "\033[1m"
	ansi_reset  = "\033[0m"
	ansi_cyan   = "\033[36m"
	ansi_yellow = "\033[33m"
	ansi_green  = "\033[32m"
	ansi_dim    = "\033[2m"
	ansi_red    = "\033[31m"
)

func (c *Create_struct) Create_run_method(args []string) error {
	flags := flagutils.Flag_parser_util(args)
	shouldStart := false
	if _, ok := flags["start"]; ok {
		shouldStart = true
	}

	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}
	client := jira.Jira_client_const(cfg)
	cache, err := cacheutils.Load_cache_util()
	if err != nil {
		return errors.Error_wrap_util("failed to load cache", err)
	}

	// 1. Draft Check
	draft := cache.TicketDraft
	var draftChoice string
	if draft != nil && time.Since(draft.UpdatedAt) < 30*time.Minute {
		fmt.Printf("\n%sFound an unsaved draft from %s%s\n", ansi_yellow, draft.UpdatedAt.Format("15:04:05"), ansi_reset)
		fmt.Printf("Project: %s | Type: %s | Summary: %s\n", draft.ProjectKey, draft.IssueType, draft.Summary)
		choice, err := promptutils.Prompt_string_util("Do you want to resume this draft? [Y/n]", "Y")
		if err == nil && strings.ToLower(strings.TrimSpace(choice)) != "n" {
			draftChoice = "y"
		} else {
			// Clear draft
			cache.TicketDraft = nil
			cacheutils.Save_cache_util(cache)
			draft = nil
		}
	} else if draft != nil {
		// Expired draft
		cache.TicketDraft = nil
		cacheutils.Save_cache_util(cache)
		draft = nil
	}

	if draft == nil {
		draft = &cacheutils.TicketDraftStruct{}
	}

	// 2. Project Selection
	if draftChoice != "y" || draft.ProjectKey == "" {
		fmt.Printf("%sFetching projects...%s\n", ansi_dim, ansi_reset)
		projects, err := client.Jira_get_projects_method()
		if err != nil {
			return errors.Error_wrap_util("failed to fetch projects", err)
		}
		
		// Order projects by MRU
		var displayProjects []string
		mruMap := make(map[string]bool)
		for _, p := range cache.RecentProjects {
			for _, jp := range projects {
				if jp.Key == p {
					displayProjects = append(displayProjects, fmt.Sprintf("%s - %s", jp.Key, jp.Name))
					mruMap[jp.Key] = true
					break
				}
			}
		}
		for _, jp := range projects {
			if !mruMap[jp.Key] {
				displayProjects = append(displayProjects, fmt.Sprintf("%s - %s", jp.Key, jp.Name))
			}
		}

		if len(displayProjects) == 0 {
			return errors.Error_wrap_util("no projects found", nil)
		}

		fmt.Printf("\n%sSelect Project:%s\n", ansi_bold, ansi_reset)
		selectedProj, ok := ttyutils.Tty_picker_util(displayProjects, "Project")
		if !ok {
			return saveDraft(cache, draft)
		}
		draft.ProjectKey = strings.Split(selectedProj, " - ")[0]
	}

	// 3. Issue Type Selection
	if draftChoice != "y" || draft.IssueTypeID == "" {
		fmt.Printf("%sFetching issue types for %s...%s\n", ansi_dim, draft.ProjectKey, ansi_reset)
		issueTypes, err := client.Jira_get_issue_types_method(draft.ProjectKey)
		if err != nil {
			return errors.Error_wrap_util("failed to fetch issue types", err)
		}

		var displayTypes []string
		for _, it := range issueTypes {
			displayTypes = append(displayTypes, it.Name)
		}

		fmt.Printf("\n%sSelect Issue Type:%s\n", ansi_bold, ansi_reset)
		selectedType, ok := ttyutils.Tty_picker_util(displayTypes, "Issue Type")
		if !ok {
			return saveDraft(cache, draft)
		}
		draft.IssueType = selectedType
		for _, it := range issueTypes {
			if it.Name == selectedType {
				draft.IssueTypeID = it.ID
				break
			}
		}
	}

	// 4. Summary
	if draftChoice != "y" || draft.Summary == "" {
		summary, err := promptutils.Prompt_string_util("Ticket Summary", draft.Summary)
		if err != nil {
			return saveDraft(cache, draft)
		}
		if strings.TrimSpace(summary) == "" {
			fmt.Println("Summary cannot be empty.")
			return saveDraft(cache, draft)
		}
		draft.Summary = strings.TrimSpace(summary)
	}

	// 5. Description
	if draftChoice != "y" || draft.Description == "" {
		desc, err := promptutils.Prompt_multiline_util("Description (Press Enter on an empty line to finish)")
		if err != nil {
			return saveDraft(cache, draft)
		}
		draft.Description = desc
	}

	// 6. Attachment (Optional)
	if draftChoice != "y" || draft.Attachment == "" {
		attachPath, err := promptutils.Prompt_string_util("Attachment File Path (Optional, leave blank to skip)", draft.Attachment)
		if err != nil {
			return saveDraft(cache, draft)
		}
		
		attachPath = strings.TrimSpace(attachPath)
		if attachPath != "" {
			info, err := os.Stat(attachPath)
			if err != nil {
				fmt.Printf("%sFile not found or inaccessible: %s%s\n", ansi_red, attachPath, ansi_reset)
				return saveDraft(cache, draft)
			}
			if info.Size() > 25*1024*1024 {
				fmt.Printf("%sFile exceeds 25MB limit. Please attach a smaller file.%s\n", ansi_red, ansi_reset)
				return saveDraft(cache, draft)
			}
			draft.Attachment = attachPath
		}
	}

	// 7. Get myself (for assignee)
	myself, err := client.Get_myself_method()
	if err != nil {
		return errors.Error_wrap_util("failed to fetch user info", err)
	}

	// 8. Confirm
	fmt.Printf("\n%s--- Review Ticket ---%s\n", ansi_bold, ansi_reset)
	fmt.Printf("Project:  %s\n", draft.ProjectKey)
	fmt.Printf("Type:     %s\n", draft.IssueType)
	fmt.Printf("Summary:  %s\n", draft.Summary)
	fmt.Printf("Assignee: Self (%s)\n", myself.DisplayName)
	if draft.Attachment != "" {
		fmt.Printf("Attach:   %s\n", draft.Attachment)
	}
	fmt.Println("---------------------")

	fmt.Printf("\nCreate this ticket? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	confirm := "n"
	if scanner.Scan() {
		confirm = strings.TrimSpace(scanner.Text())
	}
	if strings.ToLower(confirm) != "y" {
		return saveDraft(cache, draft)
	}

	// 9. Execute Creation
	fmt.Printf("\n%sCreating ticket...%s\n", ansi_dim, ansi_reset)

	fields := map[string]interface{}{
		"project": map[string]interface{}{
			"key": draft.ProjectKey,
		},
		"issuetype": map[string]interface{}{
			"id": draft.IssueTypeID,
		},
		"summary": draft.Summary,
		"assignee": map[string]interface{}{
			"id": myself.AccountID,
		},
	}

	if draft.Description != "" {
		fields["description"] = jira.Build_adf_description_util(draft.Description)
	}

	issueKey, err := client.Jira_create_issue_method(fields)
	if err != nil {
		return errors.Error_wrap_util("failed to create issue", err)
	}

	fmt.Printf("%s✅ Ticket created: %s%s\n", ansi_green, issueKey, ansi_reset)
	fmt.Printf("URL: https://%s/browse/%s\n", cfg.Jira.Domain, issueKey)

	// 10. Handle Attachment
	if draft.Attachment != "" {
		fmt.Printf("%sUploading attachment...%s\n", ansi_dim, ansi_reset)
		err := client.Jira_add_attachment_method(issueKey, draft.Attachment)
		if err != nil {
			fmt.Printf("%sFailed to upload attachment: %v%s\n", ansi_red, err, ansi_reset)
			fmt.Printf("%sYou can manually attach it via the web interface.%s\n", ansi_yellow, ansi_reset)
		} else {
			fmt.Printf("%s✅ Attachment uploaded!%s\n", ansi_green, ansi_reset)
		}
	}

	// 11. Cleanup and update MRU
	updateMRU(cache, draft.ProjectKey)
	cache.TicketDraft = nil
	cacheutils.Save_cache_util(cache)

	// 12. --start flag
	if shouldStart {
		fmt.Printf("\n%sStarting work on %s...%s\n", ansi_cyan, issueKey, ansi_reset)
		startCmd := startcmd.Start_const(c.logger)
		return startCmd.Start_run_method([]string{issueKey})
	}

	return nil
}

func saveDraft(cache *cacheutils.Cache, draft *cacheutils.TicketDraftStruct) error {
	draft.UpdatedAt = time.Now()
	cache.TicketDraft = draft
	cacheutils.Save_cache_util(cache)
	fmt.Printf("\n%sDraft saved locally.%s\n", ansi_dim, ansi_reset)
	return nil
}

func updateMRU(cache *cacheutils.Cache, projectKey string) {
	// Remove if exists
	var newMRU []string
	for _, p := range cache.RecentProjects {
		if p != projectKey {
			newMRU = append(newMRU, p)
		}
	}
	// Prepend
	cache.RecentProjects = append([]string{projectKey}, newMRU...)
	
	// Keep max 5
	if len(cache.RecentProjects) > 5 {
		cache.RecentProjects = cache.RecentProjects[:5]
	}
}
