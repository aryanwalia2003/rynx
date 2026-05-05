package initcmd

import (
	"rynx/shared/logger"
)

type Init_struct struct {
	logger *logger.Logger_struct
}

func Init_const(log *logger.Logger_struct) *Init_struct {
	return &Init_struct{
		logger: log,
	}
}
