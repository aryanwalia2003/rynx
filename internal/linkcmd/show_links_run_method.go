package linkcmd

import (
	"fmt"
	"os/exec"
	"strings"

	"rynx/shared/cacheutils"
	"rynx/shared/errors"
)

func (l *Link_struct) Show_links_run_method(args []string) error {
	cache, err := cacheutils.Load_cache_util()
	if err != nil {
		return errors.Error_wrap_util("Failed to load cache", err)
	}

	if len(cache.BranchTickets) == 0 {
		fmt.Println("No branches are currently linked to any Jira tickets.")
		return nil
	}

	out, _ := exec.Command("git", "branch", "--show-current").Output()
	currentBranch := strings.TrimSpace(string(out))

	fmt.Printf("\n%sCurrent Branch ↔ Ticket Links%s\n", ansi_bold, ansi_reset)
	fmt.Println(strings.Repeat("─", 40))

	for branch, ticket := range cache.BranchTickets {
		prefix := "  "
		if branch == currentBranch {
			prefix = "* "
			fmt.Printf("%s%s%s%s ↔ %s%s%s\n", ansi_green, prefix, branch, ansi_reset, ansi_cyan, ticket, ansi_reset)
		} else {
			fmt.Printf("%s%s ↔ %s%s%s\n", prefix, branch, ansi_cyan, ticket, ansi_reset)
		}
	}
	fmt.Println(strings.Repeat("─", 40))

	return nil
}
