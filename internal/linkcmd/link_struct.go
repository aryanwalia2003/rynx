package linkcmd

import "rynx/shared/logger"

type Link_struct struct {
	logger *logger.Logger_struct
}

func Link_const(logger *logger.Logger_struct) *Link_struct {
	return &Link_struct{logger: logger}
}
