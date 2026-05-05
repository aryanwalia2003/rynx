package authcmd

import (
	"rynx/shared/logger"
)

type Auth_struct struct {
	logger *logger.Logger_struct
}

func Auth_const(log *logger.Logger_struct) *Auth_struct {
	return &Auth_struct{
		logger: log,
	}
}
