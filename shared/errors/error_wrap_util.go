package errors

import (
	"fmt"
)

func Error_wrap_util(context string, err error) error {
	if err == nil {
		return fmt.Errorf("%s", context)
	}
	return fmt.Errorf("%s: %w", context, err)
}
