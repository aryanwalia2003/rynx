package main

import (
	"fmt"
	"os"
	"runtime"
	"rynx/internal/application"
	"rynx/shared/logger"
)

func main() {
	checkSudo()

	log := logger.Logger_const()

	app := application.Application_const(log)

	if err := app.Application_run_method(os.Args[1:]); err != nil {
		log.Logger_error_method("Application failed", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}

func checkSudo() {
	if runtime.GOOS == "windows" {
		return
	}
	if os.Geteuid() == 0 {
		fmt.Println("\033[31mError: Running rynx with sudo is blocked to prevent permission issues with your config files.\033[0m")
		fmt.Println("Please run rynx as a normal user.")
		os.Exit(1)
	}
}

