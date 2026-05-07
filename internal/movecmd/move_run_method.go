package movecmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"rynx/internal/jira"
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

func (m *Move_struct) Move_run_method(args []string) error {
	flags := flagutils.Flag_parser_util(args)

	// Filter positional args (ticket ID is optional)
	var posArgs []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
			}
			continue
		}
		posArgs = append(posArgs, arg)
	}

	// ── Resolve ticket ID ────────────────────────────────────────────────────
	ticketID := ""
	if len(posArgs) > 0 {
		ticketID = strings.ToUpper(posArgs[0])
	}

	if ticketID == "" {
		// Try to resolve from linked branch
		branchOut, err := exec.Command("git", "branch", "--show-current").Output()
		if err != nil {
			return errors.Error_wrap_util("failed to get current git branch", err)
		}
		currentBranch := strings.TrimSpace(string(branchOut))

		cache, err := cacheutils.Load_cache_util()
		if err != nil {
			return errors.Error_wrap_util("failed to load cache", err)
		}

		linked, ok := cache.BranchTickets[currentBranch]
		if !ok || linked == "" {
			return errors.Error_wrap_util(
				fmt.Sprintf("no ticket linked to branch '%s'. Provide a ticket ID or run 'dev link <TICKET-ID>' first.", currentBranch),
				nil,
			)
		}
		ticketID = linked
	}

	// Require --status flag to be present (even if empty — it's the trigger)
	if _, hasStatus := flags["status"]; !hasStatus {
		return errors.Error_wrap_util("missing --status flag. Usage: dev move [TICKET-ID] --status", nil)
	}

	// ── Load config & build Jira client ─────────────────────────────────────
	cfg, err := configutils.Config_load_merged_util()
	if err != nil {
		return err
	}
	client := jira.Jira_client_const(cfg)

	// ── Fetch available transitions ──────────────────────────────────────────
	fmt.Printf("%sFetching available transitions for %s...%s\n", ansi_dim, ticketID, ansi_reset)
	transitions, err := client.Jira_get_transitions_method(ticketID)
	if err != nil {
		return errors.Error_wrap_util("failed to fetch transitions", err)
	}
	if len(transitions) == 0 {
		fmt.Printf("%sNo transitions available for %s. Check your Jira permissions.%s\n", ansi_yellow, ticketID, ansi_reset)
		return nil
	}

	// ── Build display list for picker ────────────────────────────────────────
	transitionNames := make([]string, len(transitions))
	for i, t := range transitions {
		transitionNames[i] = t.Name
	}

	fmt.Printf("\n%s%sSelect target status for %s:%s\n", ansi_bold, ansi_cyan, ticketID, ansi_reset)
	selectedName, ok := ttyutils.Tty_picker_util(transitionNames, "Target Status")
	if !ok {
		fmt.Println("Aborted.")
		return nil
	}

	// Find the selected transition object
	var selectedTransition jira.Jira_transition_struct
	for _, t := range transitions {
		if t.Name == selectedName {
			selectedTransition = t
			break
		}
	}

	// ── Prompt for required transition screen fields ──────────────────────────
	fieldValues := jira.Jira_transition_fields_map{}
	if len(selectedTransition.Fields) > 0 {
		fmt.Printf("\n%sThis transition requires additional fields:%s\n", ansi_dim, ansi_reset)
		for fieldKey, field := range selectedTransition.Fields {
			if !field.Required {
				continue
			}
			fieldValues, err = prompt_field_util(fieldKey, field, fieldValues)
			if err != nil {
				return err
			}
		}
	}

	// ── Confirm ───────────────────────────────────────────────────────────────
	fmt.Printf("\nMove %s%s%s to %s\"%s\"%s? [y/N]: ", ansi_bold, ticketID, ansi_reset, ansi_cyan, selectedName, ansi_reset)
	scanner := bufio.NewScanner(os.Stdin)
	confirmLine := "n"
	if scanner.Scan() {
		confirmLine = strings.TrimSpace(scanner.Text())
	}
	if strings.ToLower(confirmLine) != "y" {
		fmt.Println("Aborted.")
		return nil
	}

	// ── Execute transition ────────────────────────────────────────────────────
	fmt.Printf("%sApplying transition...%s\n", ansi_dim, ansi_reset)
	if err := client.Jira_do_transition_method(ticketID, selectedTransition.ID, fieldValues); err != nil {
		return errors.Error_wrap_util("transition failed", err)
	}

	fmt.Printf("\n%s✅ %s moved to \"%s\" successfully!%s\n", ansi_green, ticketID, selectedName, ansi_reset)
	return nil
}

// prompt_field_util handles prompting for a single required transition field.
// Supported types: string/text, resolution (system), assignee (system), allowedValues (select list).
func prompt_field_util(fieldKey string, field jira.Jira_transition_field_struct, into jira.Jira_transition_fields_map) (jira.Jira_transition_fields_map, error) {
	label := field.Name
	if label == "" {
		label = fieldKey
	}

	// Select list (allowedValues present)
	if len(field.AllowedValues) > 0 {
		options := make([]string, len(field.AllowedValues))
		for i, av := range field.AllowedValues {
			options[i] = av.Name
		}
		fmt.Printf("\n%s%s:%s\n", ansi_bold, label, ansi_reset)
		selected, ok := ttyutils.Tty_picker_util(options, label)
		if !ok {
			return into, errors.Error_wrap_util("selection cancelled for required field: "+label, nil)
		}
		// Find the ID for the selected name
		for _, av := range field.AllowedValues {
			if av.Name == selected {
				into[fieldKey] = map[string]string{"id": av.ID}
				break
			}
		}
		return into, nil
	}

	// System field: assignee
	if field.Schema.System == "assignee" || fieldKey == "assignee" {
		val, err := promptutils.Prompt_string_util(fmt.Sprintf("%s (Jira account ID or email)", label), "")
		if err != nil {
			return into, err
		}
		if strings.TrimSpace(val) != "" {
			into[fieldKey] = map[string]string{"id": strings.TrimSpace(val)}
		}
		return into, nil
	}

	// Fallback: plain text / string
	val, err := promptutils.Prompt_string_util(label, "")
	if err != nil {
		return into, err
	}
	if strings.TrimSpace(val) != "" {
		into[fieldKey] = strings.TrimSpace(val)
	}
	return into, nil
}
