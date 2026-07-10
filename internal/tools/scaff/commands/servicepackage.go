package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/mitchellh/cli"
)

type ServicePackageCommand struct {
	Ui cli.Ui
}

type ServicePackageData struct {
	ServicePackage string
	ProviderName   string
	Typed          bool

	Config *helpers.Configuration
}

var _ cli.Command = ServicePackageCommand{}

func (c ServicePackageCommand) Help() string {
	return "TODO"
}

func (c ServicePackageCommand) Run(args []string) int {
	data := &ServicePackageData{
		Config: config,
	}
	if err := data.ParseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.generateRegistration(); err != nil {
		c.Ui.Error(fmt.Sprintf("creating service package: %+v", err))
		return 1
	}

	if err := data.generateClientScaffold(); err != nil {
		c.Ui.Error(fmt.Sprintf("creating service package: %+v", err))
		return 1
	}

	return 0
}

func (d *ServicePackageData) ParseArgs(args []string) (errors []error) {
	servicePackageSet := flag.NewFlagSet("servicepackage", flag.ExitOnError)
	servicePackageSet.StringVar(&d.ServicePackage, "servicepackage", "", "(Required) the name of the service package")
	servicePackageSet.BoolVar(&d.Typed, "typed", config.TypedSDK, "(Optional) Use the Typed SDK")

	err := servicePackageSet.Parse(args)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	if d.ServicePackage == "" {
		errors = append(errors, fmt.Errorf("required option `-servicepackage` missing\n"))
	}

	return
}

func (d ServicePackageData) generateRegistration() error {
	// TODO - Generate initial `registration.go`
	providerName, err := helpers.ProviderName()
	if err != nil {
		return err
	}
	d.ProviderName = *providerName

	tpl := template.Must(template.New("registration.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/registration.gotpl"))

	servicePackagePath := fmt.Sprintf("%s/%s", d.Config.ServicePackagesPath, d.ServicePackage)

	_, err = os.Stat(servicePackagePath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(servicePackagePath, 0755)
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("permission denied creating %s: %+v", servicePackagePath, err)
		}
		if err != nil {
			return fmt.Errorf("creating Service Package Path %s: %+v", servicePackagePath, err)
		}
	}

	outputPath := fmt.Sprintf("%s/registration.go", servicePackagePath)

	f, err := os.Create(outputPath)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("service package %s and registration already exists", d.ServicePackage)
	}

	err = tpl.Execute(f, d)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed writing to file: %+v", err)
	}

	// Make sure the generated template complies with the local Go's view of how it should be fmt'd
	if err := helpers.GoFmt(outputPath); err != nil {
		return err
	}

	return nil
}

func (d ServicePackageData) generateClientScaffold() error {
	providerName, err := helpers.ProviderName()
	if err != nil {
		return err
	}
	d.ProviderName = *providerName

	servicePackagePath := fmt.Sprintf("%s/%s", d.Config.ServicePackagesPath, d.ServicePackage)

	tpl := template.Must(template.New("client.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/client.gotpl"))
	clientPath := fmt.Sprintf("%s/client", servicePackagePath)

	err = os.MkdirAll(clientPath, 0755)
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("permission denied creating %q: %+v", clientPath, err)
	}
	if err != nil {
		return fmt.Errorf("creating client directory for Service Package %q: %+v", servicePackagePath, err)
	}

	outputPath := fmt.Sprintf("%s/client.go", clientPath)
	f, err := os.Create(outputPath)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("client.go already exists for service package %q: %+v", d.ServicePackage, err)
	}

	err = tpl.Execute(f, d)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed writing to file: %+v", err.Error())
	}

	return nil
}

func (c ServicePackageCommand) Synopsis() string {
	return "Creates a directory for a new Service Package and scaffolds out the basics to use it."
}
