package viewcmd

import "rynx/shared/logger"

type View_struct struct {
	logger *logger.Logger_struct
}

type View_interface interface {
	View_run_method(args []string) error
}

func View_const(log *logger.Logger_struct) View_interface {
	return &View_struct{
		logger: log,
	}
}
