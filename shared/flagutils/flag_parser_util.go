package flagutils

import (
	"strings"
)

func Flag_parser_util(args []string) map[string]string {
	flags := make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			name := strings.TrimPrefix(strings.TrimPrefix(arg, "-"), "-")
			flags[name] = "" // Mark as present
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags[name] = args[i+1]
				i++
			}
		}
	}
	return flags
}
