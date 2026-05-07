package movecmd

import "rynx/shared/logger"

type Move_struct struct {
	logger *logger.Logger_struct
}

func Move_const(logger *logger.Logger_struct) *Move_struct {
	return &Move_struct{logger: logger}
}
