package application

import (
	"rynx/internal/authcmd"
	"rynx/internal/initcmd"
	"rynx/internal/linkcmd"
	"rynx/internal/movecmd"
	"rynx/internal/prcmd"
	"rynx/internal/startcmd"
	"rynx/internal/tickets"
	"rynx/internal/viewcmd"
	"rynx/shared/errors"
)

func (a *Application_struct) Application_run_method(args []string) error {
	if len(args) == 0 {
		return errors.Error_wrap_util("no command provided", nil)
	}

	command := args[0]

	if command == "init" {
		initCmd := initcmd.Init_const(a.logger)
		return initCmd.Init_run_method()
	}

	if command == "auth" {
		authCmd := authcmd.Auth_const(a.logger)
		return authCmd.Auth_run_method()
	}

	if command == "tickets" {
		ticketsCmd := tickets.Tickets_const(a.logger)
		return ticketsCmd.Tickets_run_method(args[1:])
	}

	if command == "ping" {
		return a.Application_ping_method()
	}

	if command == "view" {
		viewCmd := viewcmd.View_const(a.logger)
		return viewCmd.View_run_method(args[1:])
	}

	if command == "pr" {
		prCmd := prcmd.PR_const(a.logger)
		return prCmd.PR_run_method(args[1:])
	}

	if command == "start" {
		startCmd := startcmd.Start_const(a.logger)
		return startCmd.Start_run_method(args[1:])
	}

	if command == "link" {
		linkCmd := linkcmd.Link_const(a.logger)
		return linkCmd.Link_run_method(args[1:])
	}

	if command == "show-links" {
		linkCmd := linkcmd.Link_const(a.logger)
		return linkCmd.Show_links_run_method(args[1:])
	}

	if command == "unlink" {
		linkCmd := linkcmd.Link_const(a.logger)
		return linkCmd.Unlink_run_method(args[1:])
	}

	if command == "move" {
		moveCmd := movecmd.Move_const(a.logger)
		return moveCmd.Move_run_method(args[1:])
	}

	return errors.Error_wrap_util("unknown command: "+command, nil)
}
