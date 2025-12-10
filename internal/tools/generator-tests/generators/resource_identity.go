// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package generators

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/templatehelpers"
	"github.com/mitchellh/cli"
)

var riOutputFileFmt = "../../services/%s/%s_resource_identity_gen_test.go"

type ResourceIdentityCommand struct {
	Ui cli.Ui
}

type resourceIdentityData struct {
	ResourceName       string
	IdentityProperties string
	PropertyNameMap    map[string]string
	ServicePackageName string
	BasicTestParams    string
	TestParams         []string
	KnownValues        string
	KnownValueMap      map[string]string
	CompareValues      string
	CompareValueMap    map[string]string
	TestName           string
}

var _ cli.Command = &ResourceIdentityCommand{}

func (c *ResourceIdentityCommand) Help() string {
	return `
Usage: generate-resource-identity [args]
Required args:
	- resource-name [string]
		the name of the resource to generate the Resource Identity test for, the 'azurerm_' prefix is not required.
	- properties [string]
		the schema exposed properties that make up the Resource Identity values. e.g. -properties "resource_group_name, site_name". ID properties that are not part of the schema, such as 'subscription_id', should not be included here.
		If the schema property name does not match the corresponding value in the ID, these should be specified together as [id_name]:[schema_name]
	- service-package-name [string]
		the name of the Service Package the resource belongs to. This forms part of the output path for the generated file.

Optional args:
	- test-params [string]
		'test-params' specifies any additional parameters that need to be passed to the 'basic' config for the resource type. e.g. '-test-params blah' == r.basic(data, "blah")
	- known-values [string]
		'known-values' specifies discriminated values that are not exposed in the schema. This is used to differentiate between resources that use the same ID type, but are discrete resources in the provider. e.g. azurerm_windows_web_app and azurerm_linux_web_app
		If the value for a 'known-value' is a CSV, replace the comma with a semi-colon to allow the parser to replace it for you. (see below for a full example)
	- compare-values [string]
		'compare-values' specifies resource identity values that do not have a one to one relationship with any values in the schema or state (i.e. the schema references a parent resource id but the resource identity includes the pieces of that parent resource id).
	- test-name [string]
		'test-name' specifies the test config name that will be used to test Resource Identity. Defaults to 'basic'.

Example:
generate-resource-identity -resource-name some_azure_resource -properties "resource_group_name,some_property" -test-params "customSku" -known-values "subscription_id:data.Subscriptions.Primary,kind:someApp;linux" -compare-values "parent_resource_name:parent_resource_id,resource_group_name:parent_resource_id"

Caveats and TODOs:
requires that the basic test for the resource is already present and has the name 'basic' for the config. TODO - Can be extended to make this configurable.
Expects that the test resource type is already declared in the test package for the service. (e.g. type LinuxFunctionAppResource struct{})
`
}

func (c *ResourceIdentityCommand) Synopsis() string {
	return "TODO - Write Synopsis for ResourceIdentityCommand"
}

func (c *ResourceIdentityCommand) Run(args []string) int {
	data := &resourceIdentityData{}

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

func (d *resourceIdentityData) parseArgs(args []string) (errors []error) {
	argSet := flag.NewFlagSet("ri", flag.ExitOnError)

	argSet.StringVar(&d.ResourceName, "resource-name", "", "(Required) the name of the resource to generate the resource identity test for.")
	argSet.StringVar(&d.IdentityProperties, "properties", "", "(Required) a comma separated list of schema property names that make up the resource identity for this resource. Do not include 'known' values here, only schema comparisons are supported.")
	argSet.StringVar(&d.ServicePackageName, "service-package-name", "", "(Required) the path to the directory containing the service package to write the generated test to.")
	argSet.StringVar(&d.BasicTestParams, "test-params", "", "(Optional) comma separated list of additional properties that need to be passed to the basic test config for this resource.")
	argSet.StringVar(&d.KnownValues, "known-values", "", "(Optional) comma separated list of known (aka discriminated) value names and their values for this resource type, formatted as [attribute_name]:[attribute value]. e.g. `kind:linux;functionapp,foo:bar`")
	argSet.StringVar(&d.CompareValues, "compare-values", "", "(Optional) comma separated list of resource identity names that are contained within a schema property value, formatted as [attribute_name]:[attribute value]. e.g. `parent_name:parent_resource_id;resource_group_name,parent_resource_id`")
	argSet.StringVar(&d.TestName, "test-name", "basic", "(Optional) the name of the config that will be used to test Resource Identity. Defaults to `basic`.")

	if err := argSet.Parse(args); err != nil {
		errors = append(errors, err)
		return
	}

	// check we have the essentials
	switch {
	case d.ResourceName == "":
		errors = append(errors, fmt.Errorf("resource name is required"))
	case d.ServicePackageName == "":
		errors = append(errors, fmt.Errorf("service-package-path is required"))
	}

	// d.PropertyNameMap = strings.Split(d.IdentityProperties, ",")
	if len(d.IdentityProperties) > 0 {
		d.PropertyNameMap = map[string]string{}
		propertiesList := strings.Split(d.IdentityProperties, ",")
		for _, property := range propertiesList {
			v := strings.Split(property, ":")
			switch len(v) {
			case 1:
				d.PropertyNameMap[v[0]] = v[0]
			case 2:
				d.PropertyNameMap[v[0]] = v[1]
			default:
				errors = append(errors, fmt.Errorf("invalid property name: %s", property))
				return
			}
		}
	}

	if len(d.BasicTestParams) > 0 {
		d.TestParams = strings.Split(d.BasicTestParams, ",")
	}

	if len(d.KnownValues) > 0 {
		d.KnownValueMap = make(map[string]string)
		kv := strings.Split(d.KnownValues, ",")
		// if len(kv)%2 != 0 {
		// 	errors = append(errors, fmt.Errorf("known-values must be a list of an even number of name/values (comma separated values should be represented with semi-colon for replacement later) e.g. 'var1:val1,var2:val2;val3'"))
		// }

		for _, v := range kv {
			vParts := strings.Split(v, ":")
			if len(vParts) != 2 {
				errors = append(errors, fmt.Errorf("invalid property format in known-values: '%s'", v))
				return
			}
			d.KnownValueMap[vParts[0]] = strings.ReplaceAll(vParts[1], ";", ",")
		}
	}

	if len(d.CompareValues) > 0 {
		d.CompareValueMap = make(map[string]string)
		kv := strings.Split(d.CompareValues, ",")

		for _, v := range kv {
			vParts := strings.Split(v, ":")
			if len(vParts) != 2 {
				errors = append(errors, fmt.Errorf("invalid property format in known-values: '%s'", v))
				return
			}
			d.CompareValueMap[vParts[0]] = strings.ReplaceAll(vParts[1], ";", ",")
		}
	}

	return
}

func (d *resourceIdentityData) exec() error {
	tpl := template.Must(template.New("identity_test.gotpl").Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, "templates/identity_test.gotpl"))

	outputPath := fmt.Sprintf(riOutputFileFmt, d.ServicePackageName, d.ResourceName)

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
		return fmt.Errorf("failed writing output test file (%s): %s", outputPath, err.Error())
	}

	if err := templatehelpers.GoImports(outputPath); err != nil {
		return err
	}

	return nil
}
