package main

import (
	"fmt"
	"os"
	"rynx/internal/application"
	"rynx/shared/logger"
)

func main() {
	log := logger.Logger_const()
	app := application.Application_const(log)

	if err := app.Application_run_method(os.Args[1:]); err != nil {
		log.Logger_error_method("Application failed", err)
		os.Exit(1)
	}
	
	fmt.Println("Done.")
}
