package startcmd

import (
	"rynx/shared/logger"
)

type Start_struct struct {
	logger *logger.Logger_struct
}

func Start_const(log *logger.Logger_struct) *Start_struct {
	return &Start_struct{
		logger: log,
	}
}
