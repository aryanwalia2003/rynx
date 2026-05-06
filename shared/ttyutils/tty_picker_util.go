package ttyutils

import (
	"fmt"
	"os"
	"rynx/shared/fuzzyutils"
	"sort"
	"strings"
)

const (
	ansi_reset      = "\033[0m"
	ansi_highlight  = "\033[33m"
	ansi_selected   = "\033[1;32m"
	ansi_dim        = "\033[2m"
	ansi_clear_line = "\033[2K\r"
	ansi_cursor_up  = "\033[A"
)

func Tty_picker_util(items []string, prompt string) (string, bool) {
	old, err := Tty_raw_enable_util()
	if err != nil {
		return "", false
	}
	defer Tty_raw_restore_util(old)

	query := ""
	cursor := 0
	buf := make([]byte, 3)
	lastLines := 0

	render := func(filtered []fuzzyutils.Fuzzy_match_struct) {
		// Move cursor up to overwrite previous render
		for i := 0; i < lastLines; i++ {
			fmt.Print(ansi_cursor_up + ansi_clear_line)
		}

		lines := 0

		// Header: prompt
		fmt.Printf("  %s%s%s\n", ansi_dim, prompt, ansi_reset)
		lines++
		// Query line
		fmt.Printf("  %s> %s%s_%s\n", ansi_highlight, ansi_reset, query, ansi_reset)
		lines++
		// Blank line
		fmt.Print("\n")
		lines++

		limit := len(filtered)
		if limit > 15 {
			limit = 15
		}
		for i := 0; i < limit; i++ {
			m := filtered[i]
			if i == cursor {
				fmt.Printf("  %s▶ %s%s\n", ansi_selected, m.Item, ansi_reset)
			} else {
				fmt.Printf("    %s\n", m.Item)
			}
			lines++
		}
		if len(filtered) > 15 {
			fmt.Printf("  %s... (%d more)%s\n", ansi_dim, len(filtered)-15, ansi_reset)
			lines++
		}

		lastLines = lines
	}

	// Initial blank render so first cursor-up works correctly
	fmt.Print("\n") // one line so the first render has something to go up to
	lastLines = 1

	for {
		filtered := fuzzyutils.Fuzzy_filter_util(items, query)
		sort.Slice(filtered, func(i, j int) bool { return filtered[i].Score > filtered[j].Score })
		if cursor >= len(filtered) {
			cursor = max_util(0, len(filtered)-1)
		}
		render(filtered)

		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}

		switch {
		case buf[0] == 13: // Enter
			if len(filtered) > 0 {
				// Clear the picker UI before returning
				for i := 0; i < lastLines; i++ {
					fmt.Print(ansi_cursor_up + ansi_clear_line)
				}
				return filtered[cursor].Item, true
			}
		case buf[0] == 27 && n == 1: // Esc
			for i := 0; i < lastLines; i++ {
				fmt.Print(ansi_cursor_up + ansi_clear_line)
			}
			return "", false
		case buf[0] == 3: // Ctrl-C
			for i := 0; i < lastLines; i++ {
				fmt.Print(ansi_cursor_up + ansi_clear_line)
			}
			return "", false
		case buf[0] == 27 && n == 3 && buf[1] == 91 && buf[2] == 65: // Up
			if cursor > 0 {
				cursor--
			}
		case buf[0] == 27 && n == 3 && buf[1] == 91 && buf[2] == 66: // Down
			cursor++
		case buf[0] == 127: // Backspace
			if len(query) > 0 {
				query = query[:len(query)-1]
				cursor = 0
			}
		default:
			if buf[0] >= 32 && buf[0] < 127 {
				query += strings.ToLower(string(buf[:n]))
				cursor = 0
			}
		}
	}
}

func max_util(a, b int) int {
	if a > b {
		return a
	}
	return b
}

