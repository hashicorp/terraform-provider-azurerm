// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/commands"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func main() {
	os.Exit(realMain(os.Args[1:]))
}

func realMain(args []string) int {
	var ui cli.Ui = &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		InfoColor:  cli.UiColorNone,

		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	Commands = map[string]cli.CommandFactory{
		"document": func() (cli.Command, error) {
			return &commands.DocumentCommand{
				Ui: ui,
			}, nil
		},
		"config": func() (cli.Command, error) {
			return &commands.GenConfigCommand{
				Ui: ui,
			}, nil
		},
		"servicepackage": func() (cli.Command, error) {
			return &commands.ServicePackageCommand{
				Ui: ui,
			}, nil
		},
		"generate": func() (cli.Command, error) {
			return &commands.GenerateCommand{
				Ui: ui,
			}, nil
		},
		"upgrade list": func() (cli.Command, error) {
			return &commands.UpgradeCommand{
				Ui: ui,
			}, nil
		},
		"upgrade typed": func() (cli.Command, error) {
			return &commands.TypedUpgradeCommand{
				Ui: ui,
			}, nil
		},
		"regen": func() (cli.Command, error) {
			return &commands.RegenCommand{
				Ui: ui,
			}, nil
		},
	}

	scaff := cli.CLI{
		Args:     args,
		Commands: Commands,
		Name:     "scaff",
		Version:  "0.1.0",
	}

	exitStatus, err := scaff.Run()
	if err != nil {
		log.Println(err)
	}
	return exitStatus
}
