package linkcmd

import (
	"fmt"
	"strings"

	"rynx/shared/cacheutils"
	"rynx/shared/errors"
)

func (l *Link_struct) Unlink_run_method(args []string) error {
	if len(args) == 0 {
		return errors.Error_wrap_util("Please provide a ticket ID to unlink (e.g. dev unlink ZZDPA-123)", nil)
	}

	ticketID := strings.ToUpper(args[0])

	cache, err := cacheutils.Load_cache_util()
	if err != nil {
		return errors.Error_wrap_util("Failed to load cache", err)
	}

	found := false
	for branch, linkedTicket := range cache.BranchTickets {
		if linkedTicket == ticketID {
			delete(cache.BranchTickets, branch)
			fmt.Printf("%s✅ Unlinked %s from branch '%s'.%s\n", ansi_green, ticketID, branch, ansi_reset)
			found = true
		}
	}

	if !found {
		fmt.Printf("%sWarning: No branches found linked to ticket %s.%s\n", ansi_yellow, ticketID, ansi_reset)
		return nil
	}

	cacheutils.Save_cache_util(cache)
	return nil
}
