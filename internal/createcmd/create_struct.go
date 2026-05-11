package createcmd

import (
	"rynx/shared/logger"
)

type Create_struct struct {
	logger *logger.Logger_struct
}

func Create_const(log *logger.Logger_struct) *Create_struct {
	return &Create_struct{
		logger: log,
	}
}
