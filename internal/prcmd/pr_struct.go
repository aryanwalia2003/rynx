package prcmd

import (
	"rynx/shared/logger"
)

type PR_struct struct {
	logger *logger.Logger_struct
}

func PR_const(log *logger.Logger_struct) *PR_struct {
	return &PR_struct{
		logger: log,
	}
}
