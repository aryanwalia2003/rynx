package startcmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"rynx/internal/jira"
	"rynx/shared/cacheutils"
	"rynx/shared/configutils"
	"rynx/shared/errors"
	"rynx/shared/promptutils"
	"rynx/shared/stringutils"
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

func (s *Start_struct) Start_run_method(args []string) error {
	if len(args) == 0 {
		return errors.Error_wrap_util("Please provide a ticket ID (e.g. dev start ZZDPA-123)", nil)
	}

	ticketID := strings.ToUpper(args[0])

	// 1. Check for dirty workspace
	statusOut, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return errors.Error_wrap_util("Failed to check git status", err)
	}
	if len(strings.TrimSpace(string(statusOut))) > 0 {
		fmt.Printf("%sWorkspace is dirty!%s\n", ansi_red, ansi_reset)
		fmt.Println("Please commit or stash your changes before starting a new ticket.")
		return nil
	}

	// 2. Fetch Jira details
	fmt.Printf("%sFetching details for %s from Jira...%s\n", ansi_dim, ticketID, ansi_reset)
	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}
	client := jira.Jira_client_const(cfg)
	detail, err := client.Jira_get_issue_detail_method(ticketID)
	if err != nil {
		return errors.Error_wrap_util("Failed to fetch Jira ticket. Make sure Jira is reachable and the ticket exists.", err)
	}

	fmt.Printf("\n%s%s[%s]%s %s\n\n", ansi_bold, ansi_cyan, detail.Key, ansi_reset, detail.Summary)

	// 3. Ask for base branch and checkout
	baseBranch, err := promptutils.Prompt_string_util("Base branch to cut from", "main")
	if err != nil {
		return err
	}

	fmt.Printf("%sChecking out %s and pulling latest changes...%s\n", ansi_dim, baseBranch, ansi_reset)
	checkoutCmd := exec.Command("git", "checkout", baseBranch)
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr
	if err := checkoutCmd.Run(); err != nil {
		return errors.Error_wrap_util("Failed to checkout base branch", err)
	}

	pullCmd := exec.Command("git", "pull", "origin", baseBranch)
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return errors.Error_wrap_util("Failed to pull latest changes", err)
	}

	// 4. Prefix and slug generation
	prefix, err := promptutils.Prompt_string_util("Branch prefix (e.g. feat, bug, fix)", "feat")
	if err != nil {
		return err
	}

	slug := stringutils.Slugify_util(detail.Summary)
	proposedBranch := fmt.Sprintf("%s/%s-%s", prefix, ticketID, slug)

	// 5. Propose and confirm
	fmt.Printf("\n%sProposed branch name:%s\n", ansi_bold, ansi_reset)
	finalBranch, err := promptutils.Prompt_string_util("Edit or confirm branch name", proposedBranch)
	if err != nil {
		return err
	}

	// 6. Check if branch exists
	branchExistsOut, _ := exec.Command("git", "rev-parse", "--verify", finalBranch).CombinedOutput()
	if len(branchExistsOut) > 0 && !strings.Contains(string(branchExistsOut), "fatal") {
		fmt.Printf("\n%sBranch %s already exists!%s\n", ansi_yellow, finalBranch, ansi_reset)
		switchPrompt, err := promptutils.Prompt_string_util("Switch to it? [Y/n]", "Y")
		if err != nil {
			return err
		}
		if strings.ToLower(strings.TrimSpace(switchPrompt)) != "n" {
			exec.Command("git", "checkout", finalBranch).Run()
			fmt.Printf("%sSwitched to %s.%s\n", ansi_green, finalBranch, ansi_reset)
		} else {
			return nil
		}
	} else {
		// 7. Create branch
		fmt.Printf("\n%sCreating branch %s...%s\n", ansi_dim, finalBranch, ansi_reset)
		createCmd := exec.Command("git", "checkout", "-b", finalBranch)
		createCmd.Stdout = os.Stdout
		createCmd.Stderr = os.Stderr
		if err := createCmd.Run(); err != nil {
			return errors.Error_wrap_util("Failed to create branch", err)
		}
		fmt.Printf("%s✅ Branch created successfully!%s\n", ansi_green, ansi_reset)
	}

	// 8. Cache mapping
	cache, err := cacheutils.Load_cache_util()
	if err == nil {
		cache.BranchTickets[finalBranch] = ticketID
		cacheutils.Save_cache_util(cache)
	}

	return nil
}
