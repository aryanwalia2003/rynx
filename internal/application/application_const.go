package application

import (
	"rynx/shared/logger"
)

func Application_const(log *logger.Logger_struct) *Application_struct {
	return &Application_struct{
		logger: log,
	}
}
