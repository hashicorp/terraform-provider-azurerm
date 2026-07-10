package commands

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/mitchellh/cli"
)

type GenConfigCommand struct {
	Ui cli.Ui
}

type localConfig helpers.Configuration

var _ cli.Command = GenConfigCommand{}

func (c GenConfigCommand) Run(_ []string) int {
	l := localConfig(*config)
	if err := l.generate(); err != nil {
		c.Ui.Error(fmt.Sprintf("generating config file: %+v", err))
	}

	return 0
}

func (c GenConfigCommand) Synopsis() string {
	return "Generate a local config file for common options."
}

func (c GenConfigCommand) Help() string {
	return `
Usage: scaff config

Generates a local configuration file for scaff which can be used to set defaults for options to some commands.
`
}

func (d *localConfig) generate() error {
	tpl := template.Must(template.New("config.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/config.gotpl"))
	outputPath := fmt.Sprintf("./%s", helpers.ConfigFileName)

	f, err := os.Create(outputPath)
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("permission denied writing config file: %+v", err)
	}
	err = tpl.Execute(f, d)
	if err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}
