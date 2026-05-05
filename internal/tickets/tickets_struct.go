package tickets

import (
	"rynx/shared/logger"
)

type Tickets_struct struct {
	logger *logger.Logger_struct
}

func Tickets_const(log *logger.Logger_struct) *Tickets_struct {
	return &Tickets_struct{
		logger: log,
	}
}
