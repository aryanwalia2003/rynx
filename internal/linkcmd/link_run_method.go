package linkcmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"rynx/internal/jira"
	"rynx/shared/cacheutils"
	"rynx/shared/configutils"
	"rynx/shared/errors"
	"rynx/shared/flagutils"
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

func (l *Link_struct) Link_run_method(args []string) error {
	flags := flagutils.Flag_parser_util(args)
	_, force := flags["force"]

	// Filter out flags from args to get positional arguments
	var posArgs []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// Skip the value if there is one (Flag_parser_util logic)
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
			}
			continue
		}
		posArgs = append(posArgs, arg)
	}

	if len(posArgs) == 0 {
		return errors.Error_wrap_util("Please provide a ticket ID (e.g. dev link ZZDPA-123)", nil)
	}

	ticketID := strings.ToUpper(posArgs[0])

	targetBranch := ""
	if len(posArgs) > 1 {
		targetBranch = posArgs[1]
	} else {
		// Get current branch
		out, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			return errors.Error_wrap_util("Failed to get current git branch", err)
		}
		targetBranch = strings.TrimSpace(string(out))
		if targetBranch == "" {
			return errors.Error_wrap_util("Not currently on any branch", nil)
		}
	}

	cache, err := cacheutils.Load_cache_util()
	if err != nil {
		return errors.Error_wrap_util("Failed to load cache", err)
	}

	// Enforce 1:1 Mapping
	for branch, linkedTicket := range cache.BranchTickets {
		if linkedTicket == ticketID && branch != targetBranch {
			if !force {
				fmt.Printf("%sTicket %s is already linked to branch '%s'.%s\n", ansi_red, ticketID, branch, ansi_reset)
				fmt.Println("Use --force to override and steal the ticket link.")
				return nil
			} else {
				fmt.Printf("%sStealing link for %s from branch '%s'...%s\n", ansi_yellow, ticketID, branch, ansi_reset)
				delete(cache.BranchTickets, branch)
			}
		}
	}

	// Prompt for rename
	fmt.Printf("%sAre you sure you want to link or do you want to name branch according to the jira ticket?%s\n", ansi_bold, ansi_reset)
	renameChoice, err := promptutils.Prompt_string_util("Y for Rename, N for continue with current branch name [y/N]", "N")
	if err != nil {
		return err
	}

	if strings.ToLower(strings.TrimSpace(renameChoice)) == "y" {
		// Flow Y (Rename)
		fmt.Printf("%sFetching Jira details for %s...%s\n", ansi_dim, ticketID, ansi_reset)
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

		prefix, err := promptutils.Prompt_string_util("Branch prefix (e.g. feat, bug, fix)", "feat")
		if err != nil {
			return err
		}

		slug := stringutils.Slugify_util(detail.Summary)
		proposedBranch := fmt.Sprintf("%s/%s-%s", prefix, ticketID, slug)

		fmt.Printf("\n%sProposed new branch name:%s\n", ansi_bold, ansi_reset)
		finalBranch, err := promptutils.Prompt_string_util("Edit or confirm branch name", proposedBranch)
		if err != nil {
			return err
		}

		fmt.Printf("\n%sRenaming branch '%s' to '%s'...%s\n", ansi_dim, targetBranch, finalBranch, ansi_reset)
		renameCmd := exec.Command("git", "branch", "-m", targetBranch, finalBranch)
		renameCmd.Stdout = os.Stdout
		renameCmd.Stderr = os.Stderr
		if err := renameCmd.Run(); err != nil {
			return errors.Error_wrap_util("Failed to rename branch", err)
		}

		// Switch to new branch
		checkoutCmd := exec.Command("git", "checkout", finalBranch)
		checkoutCmd.Stdout = os.Stdout
		checkoutCmd.Stderr = os.Stderr
		if err := checkoutCmd.Run(); err != nil {
			return errors.Error_wrap_util("Failed to checkout new branch", err)
		}

		targetBranch = finalBranch
		fmt.Printf("%s✅ Branch renamed and linked! Switched to %s.%s\n", ansi_green, finalBranch, ansi_reset)

	} else {
		// Flow N (Default/Link Only)
		// Warn if ticketID doesn't look like a Jira issue (e.g. PROJ-123)
		match, _ := regexp.MatchString(`^[A-Z]+-\d+$`, ticketID)
		if !match {
			fmt.Printf("%sWarning: '%s' does not look like a standard Jira ticket ID (e.g., PROJ-123).%s\n", ansi_yellow, ticketID, ansi_reset)
		}
		fmt.Printf("%s✅ Linked branch '%s' to %s.%s\n", ansi_green, targetBranch, ticketID, ansi_reset)
	}

	// Clean up old link for this branch if exists
	if oldTicket, exists := cache.BranchTickets[targetBranch]; exists {
		if oldTicket != ticketID {
			fmt.Printf("%sReplaced old link %s -> %s%s\n", ansi_dim, oldTicket, ticketID, ansi_reset)
		}
	}

	cache.BranchTickets[targetBranch] = ticketID
	if err := cacheutils.Save_cache_util(cache); err != nil {
		return errors.Error_wrap_util("Failed to save cache", err)
	}

	return nil
}
