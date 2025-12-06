// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commands

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/templatehelpers"
	"github.com/mitchellh/cli"
)

type ResourceCommand struct {
	Ui cli.Ui
}

var _ cli.Command = &ResourceCommand{}

type resourceData struct {
	Name                     string   `json:"name"                   hcl:"name"`
	ServicePackageName       string   `json:"service_package_name"   hcl:"service_package_name"`
	RPName                   string   `json:"resource_provider_name" hcl:"resource_provider_name"`
	ClientName               string   `json:"client_name"            hcl:"client_name"`
	Updatable                bool     `json:"updatable"              hcl:"updatable"`
	UsesLROCRUD              bool     `json:"uses_lro_crud"          hcl:"uses_lro_crud"`
	UseCreateOptions         bool     `json:"use_create_options"     hcl:"use_create_options"`
	UseReadOptions           bool     `json:"use_read_options"       hcl:"use_read_options"`
	UseUpdateOptions         bool     `json:"use_update_options"     hcl:"use_update_options"`
	UseDeleteOptions         bool     `json:"use_delete_options"     hcl:"use_delete_options"`
	ConfigValidators         bool     `json:"config_validators"      hcl:"config_validators"`
	APIVersion               string   `json:"api_version"            hcl:"api_version"`
	NoResourceGroup          bool     `json:"no_resource_group"      hcl:"no_resource_group"`
	IdType                   string   `json:"id_type"                hcl:"id_type"`
	IdTypeParts              []string `json:"id_type_parts"          hcl:"id_type_parts"`
	SDKName                  string   `json:"sdk_name"               hcl:"sdk_name"`
	IdSegments               string   `json:"id_segments"            hcl:"id_segments"`
	IDSegments               []string
	ResourceIdentitySegments []string
	SegmentCount             int
}

var outputResourceFileFmt = relativePathToRoot() + "internal/services/%s/%s_resource.go"

var outputModelFileFmt = relativePathToRoot() + "internal/services/%s/%s_resource_models.go"

func (c ResourceCommand) Run(args []string) int {
	data := &resourceData{}

	if err := data.parseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}

		return 1
	}

	if err := data.exec(); err != nil {
		c.Ui.Error(err.Error())

		log.Println(err)
		return 2
	}

	return 0
}

func (c ResourceCommand) Synopsis() string {
	return "create boilerplate for AzureRM/AD Framework Implemented resource to speed development"
}

func (c ResourceCommand) Help() string {
	return `
Usage: scaff resource -name "some_resource_name" -service_package_name="someservice" -rp_name="sql" -client_name="SomeClient" [-updatable=true] [-no_resource_group=true] -api_version="2023-08-01-preview" -id_type="commonids.SqlDatabaseId" [-sdk_name="databases"] -id_segments="SubscriptionId,ResourceGroupName,ServerName,DatabaseName" [-uses_lro_crud=true] [-use_create_options=true] [-use_read_options=true] [-use_update_options=true] [-use_delete_options=true]

Parameters:
	-name (Required) the name of the resource to scaffold the resource for.
	-service_package_name (Required) the name of the service package to scaffold the resource into.
	-rp_name (Required) the name of the resource provider of the new resource."
	-api_version (Required) the API version of the resource to scaffold. e.g. 2025-01-01")
	-id_type (Required) the type of resource to scaffold. e.g. 'commonids.AppServiceId', or 'virtualmachines.VirtualMachineId'.
	-id_segments (Required) The User-Specified Segment names for the ID, Order matters. Future versions of this command will discover this from the id_type value, probably...

	-updatable (Optional) whether the new resource can be updated. i.e. any schema property is not going to be 'ForceNew'.
	-uses_lro_crud (Optional) the new resource uses LROs for Create, Update, and Delete.
	-use_create_options (Optional) the new resource uses OperationOptions for Create.
	-use_read_options (Optional) the new resource uses OperationOptions for Read.
	-use_update_options (Optional) the new resource uses OperationOptions for Update.
	-use_delete_options (Optional) the new resource uses OperationOptions for Delete.
	-config_validators (Optional) does the resource have configuration validators.
	-no_resource_group (Optional) Set to true if the resource is not created in a resource group, or if the RG is inferred from a parent resource ID.
	-sdk_name (Optional) the name of the SDK used to manage the new resource. If omitted, the first slug of the id_type value will be used.

Example:
scaff resource -name="fw_mssql_database" -service_package_name="MSSQL" -rp_name="sql" -client_name="DatabasesClient" -updatable=true -no_resource_group=true -api_version="2023-08-01-preview" -id_type="commonids.SqlDatabaseId" -sdk_name="databases" -id_segments="SubscriptionId,ResourceGroupName,ServerName,DatabaseName" -uses_lro_crud=true -use_read_options=true

or 

go run internal/tools/scaff/main.go resource -name="fw_mssql_database" -service_package_name="MSSQL" -rp_name="sql" -client_name="DatabasesClient" -updatable=true -no_resource_group=true -api_version="2023-08-01-preview" -id_type="commonids.SqlDatabaseId" -sdk_name="databases" -id_segments="SubscriptionId,ResourceGroupName,ServerName,DatabaseName" -uses_lro_crud=true -use_read_options=true
`
}

func (d *resourceData) parseArgs(args []string) (errs []error) {
	argSet := flag.NewFlagSet("resource", flag.ExitOnError)

	argSet.StringVar(&d.Name, "name", "", "(Required) the name of the resource to scaffold the resource for.")
	argSet.StringVar(&d.ServicePackageName, "service_package_name", "", "(Required) the name of the service package to scaffold the resource into.")
	argSet.StringVar(&d.RPName, "rp_name", "", "(Required) the name of the resource provider of the new resource.")
	argSet.StringVar(&d.ClientName, "client_name", "", "(Required) the name of the client used to manage the new resource.")
	argSet.StringVar(&d.APIVersion, "api_version", "", "(Required) the API version of the resource to scaffold. e.g. 2025-01-01")
	argSet.StringVar(&d.IdType, "id_type", "", "(Required) the type of resource to scaffold. e.g. `commonids.AppServiceId`, or `virtualmachines.VirtualMachineId`.")
	argSet.StringVar(&d.IdSegments, "id_segments", "", "(Required) The User-Specified Segment names for the ID, Order matters. Future versions of this command will discover this from the id_type value, I hope...")
	argSet.BoolVar(&d.Updatable, "updatable", false, "(Optional) whether the new resource can be updated. i.e. any schema property is not going to be `ForceNew`.")
	argSet.BoolVar(&d.UsesLROCRUD, "uses_lro_crud", false, "(Optional) the new resource uses LROs for Create, Update, and Delete.")
	argSet.BoolVar(&d.UseCreateOptions, "use_create_options", false, "(Optional) the new resource uses OperationOptions for Create.")
	argSet.BoolVar(&d.UseReadOptions, "use_read_options", false, "(Optional) the new resource uses OperationOptions for Read.")
	argSet.BoolVar(&d.UseUpdateOptions, "use_update_options", false, "(Optional) the new resource uses OperationOptions for Update.")
	argSet.BoolVar(&d.UseDeleteOptions, "use_delete_options", false, "(Optional) the new resource uses OperationOptions for Delete.")
	argSet.BoolVar(&d.ConfigValidators, "config_validators", false, "(Optional) does the resource have configuration validators.")
	argSet.BoolVar(&d.NoResourceGroup, "no_resource_group", false, "(Optional) Set to true if the resource is not created in a resource group, or if the RG is inferred from a parent resource ID.")
	argSet.StringVar(&d.SDKName, "sdk_name", "", "(Optional) the name of the SDK used to manage the new resource. If omitted, the first slug of the id_type value will be used.")
	if err := argSet.Parse(args); err != nil {
		errs = append(errs, err)
		return errs
	}

	switch {
	case d.Name == "":
		errs = append(errs, errors.New("resource name is required"))
	case d.ServicePackageName == "":
		errs = append(errs, errors.New("service package name is required"))
	case d.RPName == "":
		errs = append(errs, errors.New("resource name is required"))
	case d.ClientName == "":
		errs = append(errs, errors.New("client name is required"))
	case d.APIVersion == "":
		errs = append(errs, errors.New("api version is required"))
	case d.IdType == "":
		errs = append(errs, errors.New("id_type is required"))
	case d.IdSegments == "":
		errs = append(errs, errors.New("id_segments is required"))
	}

	d.IdTypeParts = strings.Split(d.IdType, ".")
	if l := len(d.IdTypeParts); l != 2 {
		errs = append(errs, fmt.Errorf("id_type has incorrect number of segments, expected 2 got %d", l))
	}

	d.IDSegments = strings.Split(d.IdSegments, ",")
	d.SegmentCount = len(d.IDSegments) - 1

	for idx, segment := range d.IDSegments {
		switch {
		case idx == len(d.IDSegments)-1:
			d.ResourceIdentitySegments = append(d.ResourceIdentitySegments, "Name")
		default:
			d.ResourceIdentitySegments = append(d.ResourceIdentitySegments, segment)
		}
	}

	return errs
}

func (d *resourceData) exec() error {
	// Generate Resource
	tpl := template.Must(template.New("resource.gotpl").Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, "templates/resource.gotpl"))

	outputPath := fmt.Sprintf(outputResourceFileFmt, d.ServicePackageName, d.Name)

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed opening output resource file for writing: %+v", err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("failed closing output resource file for writing:", err.Error())
			os.Exit(3)
		}
	}(f)

	if err := tpl.Execute(f, d); err != nil {
		return fmt.Errorf("failed writing output resource file (%s): %s", outputPath, err.Error())
	}

	if err := templatehelpers.GoImports(outputPath); err != nil {
		return err
	}

	// Generate Resource Model(s)
	tpl = template.Must(template.New("resource_models.gotpl").Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, "templates/resource_models.gotpl"))

	outputPath = fmt.Sprintf(outputModelFileFmt, d.ServicePackageName, d.Name)

	f, err = os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed opening output models file for writing: %+v", err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("failed closing output models file for writing:", err.Error())
			os.Exit(3)
		}
	}(f)

	if err := tpl.Execute(f, d); err != nil {
		return fmt.Errorf("failed writing output models file (%s): %s", outputPath, err.Error())
	}

	if err := templatehelpers.GoImports(outputPath); err != nil {
		return err
	}

	return nil
}

func relativePathToRoot() string {
	here, err := exec.Command("git", "rev-parse", "--show-cdup").Output()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(here))
}
