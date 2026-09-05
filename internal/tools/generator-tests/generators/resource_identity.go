// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package generators

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/templatehelpers"
	"github.com/mitchellh/cli"
)

var (
	cwd, _          = os.Getwd()
	riOutputFileFmt = "/%s_resource_identity_gen_test.go"
)

type ResourceIdentityCommand struct {
	Ui cli.Ui
}

type resourceIdentityData struct {
	ResourceName           string
	IdentityProperties     string
	ResourceID             string
	ParentID               string
	PropertyNameMap        map[string]string
	ServicePackageName     string
	BasicTestParams        string
	TestParams             []string
	KnownValues            string
	KnownValueMap          map[string]string
	CompareValues          string
	CompareValueMap        map[string]string
	NoSubscriptionID       bool
	TestName               string
	TestEnvVarsFlag        string
	TestEnvVars            []string
	TestExpectNonEmptyPlan bool
	TestSequential         bool
	RewriteTag             bool
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

Optional args:
	- service-package-name [string]
		the name of the Service Package the resource belongs to. This forms part of the output path for the generated file. For Go Generate this will be picked up from the current working directory.
	- test-params [string]
		'test-params' specifies any additional parameters that need to be passed to the 'basic' config for the resource type. e.g. '-test-params blah' == r.basic(data, "blah")
	- known-values [string]
		'known-values' specifies discriminated values that are not exposed in the schema. This is used to differentiate between resources that use the same ID type, but are discrete resources in the provider. e.g. azurerm_windows_web_app and azurerm_linux_web_app
		If the value for a 'known-value' is a CSV, replace the comma with a semi-colon to allow the parser to replace it for you. (see below for a full example)
	- compare-values [string]
		'compare-values' specifies resource identity values that do not have a one to one relationship with any values in the schema or state (i.e. the schema references a parent resource id but the resource identity includes the pieces of that parent resource id).
	- no-subscription-id [bool]
		- 'no-subscription-id' can be passed to omit the default subscription ID known value in the test if the resource's ID doesn't contain subscription ID. Defaults to 'false'.
	- test-name [string]
		'test-name' specifies the test config name that will be used to test Resource Identity. Defaults to 'basic'.
	- test-env-vars [string]
		- 'test-env-vars' specifies any environment variables that are required for the test to run, the test is skipped if specified env vars are not defined.
	- test-expect-non-empty [bool]
		'test-expect-non-empty' indicates whether the test should expect (and ignore) a non-empty plan, to be used when the API does not return certain values during import.
	- test-sequential [bool]
		'test-sequential' generates a lowercase test function name (testAcc... instead of TestAcc...) and uses resource.Test instead of resource.ParallelTest.
		This is used for resources that need to run in a sequential test suite.

Example:
generate-resource-identity -resource-name some_azure_resource -properties "resource_group_name,some_property" -test-params "customSku" -known-values "kind:someApp;linux" -compare-values "parent_resource_name:parent_resource_id,resource_group_name:parent_resource_id"

Caveats and TODOs:
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

	argSet.StringVar(&d.ResourceName, "resource-name", "", "(Optional) the name of the resource to generate the resource identity test for. Defaults to parsing $GOFILE.")
	argSet.StringVar(&d.IdentityProperties, "properties", "", "(Optional) a comma separated list of schema property names that make up the resource identity for this resource. Do not include 'known' values here, only schema comparisons are supported.")
	argSet.StringVar(&d.ResourceID, "id", "", "(Optional) the Azure Resource ID string with placeholders in curly braces, e.g. /subscriptions/{subscription_id}/resourceGroups/{resource_group_name}/providers/Microsoft.Web/sites/{name}. This is a simpler alternative to -properties.")
	argSet.StringVar(&d.ParentID, "parent-id", "", "(Optional) for sub-resources with a virtual identity, specify the schema property name of the parent ID (e.g. `workspace_id`). This automatically expands into compare-values.")
	argSet.StringVar(&d.ServicePackageName, "service-package-name", "", "(Optional) the path to the directory containing the service package to write the generated test to. For Go generate this will be picked up from $GOPACKAGE or current working directory.")
	argSet.StringVar(&d.BasicTestParams, "test-params", "", "(Optional) comma separated list of additional properties that need to be passed to the basic test config for this resource.")
	argSet.StringVar(&d.TestEnvVarsFlag, "test-env-vars", "", "(Optional) comma separated list of env vars that need to be passed to the basic test config for this resource, if any of the provided env vars are missing the test is skipped.")
	argSet.StringVar(&d.KnownValues, "known-values", "", "(Optional) comma separated list of known (aka discriminated) value names and their values for this resource type, formatted as [attribute_name]:[attribute value]. e.g. `kind:linux;functionapp,foo:bar`")
	argSet.StringVar(&d.CompareValues, "compare-values", "", "(Optional) comma separated list of resource identity names that are contained within a schema property value, formatted as [attribute_name]:[attribute value]. e.g. `parent_name:parent_resource_id;resource_group_name,parent_resource_id`")
	argSet.BoolVar(&d.NoSubscriptionID, "no-subscription-id", false, "(Optional) 'no-subscription-id' can be passed to omit the default subscription ID known value in the test if the resource's ID doesn't contain subscription ID. Defaults to 'false'.")
	argSet.StringVar(&d.TestName, "test-name", "basic", "(Optional) the name of the config that will be used to test Resource Identity. Defaults to `basic`.")
	argSet.BoolVar(&d.TestExpectNonEmptyPlan, "test-expect-non-empty", false, "(Optional) Whether to expect (and ignore) a non-empty plan, to be used when the API does not return certain values during import. Defaults to `false`.")
	argSet.BoolVar(&d.TestSequential, "test-sequential", false, "(Optional) generates a lowercase test function name for use in sequential test suites. Defaults to `false`.")

	if err := argSet.Parse(args); err != nil {
		errors = append(errors, err)
		return
	}

	if d.ResourceName == "" {
		if goFile := os.Getenv("GOFILE"); goFile != "" {
			d.ResourceName = strings.TrimSuffix(goFile, "_resource.go")
		}
	}

	// check we have the essentials
	if d.ResourceName == "" {
		errors = append(errors, fmt.Errorf("`resource-name` is required (or must be run via go:generate to infer from $GOFILE)"))
	}

	// remove accidental provider prefix
	d.ResourceName = strings.TrimPrefix(d.ResourceName, "azurerm_")

	if d.ResourceID != "" && d.IdentityProperties == "" {
		d.PropertyNameMap = map[string]string{}

		// Extract placeholders like {name} or {subscription_id}
		re := regexp.MustCompile(`\{([^}]+)\}`)
		matches := re.FindAllStringSubmatch(d.ResourceID, -1)

		hasSubscriptionId := false
		for _, match := range matches {
			if len(match) > 1 {
				propName := match[1]
				// Handle property overrides like {group_id:name}
				parts := strings.Split(propName, ":")
				if len(parts) == 2 {
					d.PropertyNameMap[parts[0]] = parts[1]
				} else {
					if propName == "subscription_id" {
						hasSubscriptionId = true
					} else {
						d.PropertyNameMap[propName] = propName
					}
				}
			}
		}

		if hasSubscriptionId {
			if d.KnownValues != "" {
				d.KnownValues += ",subscription_id:data.Subscriptions.Primary"
			} else {
				d.KnownValues = "subscription_id:data.Subscriptions.Primary"
			}
		}
	} else if len(d.IdentityProperties) > 0 {
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

	if len(d.CompareValues) > 0 {
		d.CompareValueMap = make(map[string]string)
		kv := strings.Split(d.CompareValues, ",")

		for _, v := range kv {
			vParts := strings.Split(v, ":")
			if len(vParts) != 2 {
				errors = append(errors, fmt.Errorf("invalid property format in compare-values: '%s'", v))
				return
			}

			name := vParts[0]
			if name == "subscription_id" {
				// prevent duplicate `subscription_id` check
				d.NoSubscriptionID = true
			}
			d.CompareValueMap[vParts[0]] = strings.ReplaceAll(vParts[1], ";", ",")
		}
	}

	// AST Inference & Self-Healing
	d.RewriteTag = false
	goFile := os.Getenv("GOFILE")
	if goFile != "" {
		inferred, err := InferIdentityProperties(goFile)
		if err == nil {
			inferredMap := map[string]string{}
			for _, prop := range inferred.Properties {
				inferredMap[prop] = prop
			}

			if !inferred.IsVirtual && len(d.CompareValueMap) > 0 && len(d.PropertyNameMap) == 0 {
				allSameValue := true
				var commonValue string
				for k, v := range d.CompareValueMap {
					if _, ok := inferredMap[k]; ok || k == "subscription_id" {
						if commonValue == "" {
							commonValue = v
						} else if commonValue != v {
							allSameValue = false
							break
						}
					}
				}
				if allSameValue && commonValue != "" {
					inferred.IsVirtual = true
				}
			}

			if d.ParentID != "" {
				inferred.IsVirtual = true
			}

			if inferred.IsVirtual {
				if d.CompareValueMap == nil {
					d.CompareValueMap = make(map[string]string)
				}

				// If they provided -parent-id, auto-expand it!
				if d.ParentID != "" {
					for _, prop := range inferred.Properties {
						d.CompareValueMap[prop] = d.ParentID
					}
					if inferred.HasSubscriptionID {
						d.CompareValueMap["subscription_id"] = d.ParentID
						d.NoSubscriptionID = true
					}
				}

				// If they provided -compare-values (legacy), see if it perfectly matches a single parent ID
				if len(d.CompareValueMap) > 0 {
					allSameValue := true
					var commonValue string
					for k, v := range d.CompareValueMap {
						// Only check properties that belong to the identity
						if _, ok := inferredMap[k]; ok || k == "subscription_id" {
							if commonValue == "" {
								commonValue = v
							} else if commonValue != v {
								allSameValue = false
								break
							}
						}
					}

					// We can safely rewrite -compare-values "..." to -parent-id "commonValue"
					if allSameValue && commonValue != "" && d.ParentID == "" {
						d.RewriteTag = true
						d.ParentID = commonValue // act as if it was provided

						// Re-expand everything just to be safe
						for _, prop := range inferred.Properties {
							d.CompareValueMap[prop] = d.ParentID
						}
						if inferred.HasSubscriptionID {
							d.CompareValueMap["subscription_id"] = d.ParentID
						}
					}
				} else if d.ParentID != "" {
					// they already used -parent-id, we can still rewrite to drop -compare-values if it exists
					d.RewriteTag = true
				}

				// Ensure PropertyNameMap is empty so generator uses compare-values exclusively
				d.PropertyNameMap = map[string]string{}
			} else {
				if len(d.PropertyNameMap) == 0 {
					// We had no properties/id provided, use inferred
					d.PropertyNameMap = inferredMap
					if inferred.HasSubscriptionID {
						if d.KnownValues != "" {
							d.KnownValues += ",subscription_id:data.Subscriptions.Primary"
						} else {
							d.KnownValues = "subscription_id:data.Subscriptions.Primary"
						}
					}
					// We don't need to rewrite if it's already clean, but we could
				} else {
					// We had properties/id provided. Do they match inferred?
					match := true
					if len(d.PropertyNameMap) != len(inferredMap) {
						match = false
					} else {
						for k, v := range d.PropertyNameMap {
							if infV, ok := inferredMap[k]; !ok || infV != v {
								match = false
								break
							}
						}
					}

					if match {
						// We can safely drop -properties and -id from the go:generate tag
						d.RewriteTag = true
					}
				}
			}
		} else if len(d.PropertyNameMap) == 0 && d.ParentID == "" && d.CompareValues == "" {
			errors = append(errors, fmt.Errorf("neither -properties, -id, -parent-id, -compare-values nor $GOFILE AST inference succeeded: %w", err))
			return
		}
	} else if len(d.PropertyNameMap) == 0 && d.ParentID == "" && d.CompareValues == "" {
		errors = append(errors, fmt.Errorf("neither -properties, -id, -parent-id, -compare-values nor $GOFILE were provided to determine identity properties"))
		return
	}

	if len(d.BasicTestParams) > 0 {
		d.TestParams = strings.Split(d.BasicTestParams, ",")
	}

	if len(d.TestEnvVarsFlag) > 0 {
		d.TestEnvVars = strings.Split(d.TestEnvVarsFlag, ",")
	}

	if len(d.KnownValues) > 0 {
		d.KnownValueMap = make(map[string]string)
		kv := strings.Split(d.KnownValues, ",")

		for _, v := range kv {
			vParts := strings.Split(v, ":")
			if len(vParts) != 2 {
				errors = append(errors, fmt.Errorf("invalid property format in known-values: '%s'", v))
				return
			}

			name := vParts[0]
			if name == "subscription_id" {
				// prevent duplicate `subscription_id` check
				d.NoSubscriptionID = true
			}
			d.KnownValueMap[name] = strings.ReplaceAll(vParts[1], ";", ",")
		}
	}

	return
}

func (d *resourceIdentityData) exec() error {
	tpl := template.Must(template.New("identity_test.gotpl").Funcs(templatehelpers.TplFuncMap).ParseFS(Templatedir, "templates/identity_test.gotpl"))

	outputPath := cwd + fmt.Sprintf(riOutputFileFmt, d.ResourceName)
	cwdParts := strings.Split(cwd, "internal"+string(os.PathSeparator)+"services"+string(os.PathSeparator))

	// Allow service package name override if needed (unlikely)
	if d.ServicePackageName == "" {
		d.ServicePackageName = cwdParts[len(cwdParts)-1]
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed opening output resource file for writing: %+v", err.Error())
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("failed closing output resource file for writing:", err.Error())
			os.Exit(3)
		}
	}(f)

	if err := tpl.Execute(f, d); err != nil {
		return fmt.Errorf("failed writing output test file (%s): %s", outputPath, err.Error())
	}

	if err := templatehelpers.GoImports(outputPath); err != nil {
		return fmt.Errorf("failed to run goimports: %w", err)
	}

	if d.RewriteTag {
		goFile := os.Getenv("GOFILE")
		if goFile != "" {
			content, err := os.ReadFile(goFile)
			if err != nil {
				return fmt.Errorf("failed reading %s for self-healing: %w", goFile, err)
			}

			// Remove -properties "...", -id "...", and redundant flags
			reProp := regexp.MustCompile(`[ \t]*-properties[ \t]+"[^"]+"`)
			reId := regexp.MustCompile(`[ \t]*-id[ \t]+"[^"]+"`)
			reResName := regexp.MustCompile(`[ \t]*-resource-name[ \t]+[^ \t\n]+`)
			rePkgName := regexp.MustCompile(`[ \t]*-service-package-name[ \t]+[^ \t\n]+`)
			reKnownSub := regexp.MustCompile(`[ \t]*-known-values[ \t]+"subscription_id:data\.Subscriptions\.Primary"`)

			newContent := reProp.ReplaceAllString(string(content), "")
			newContent = reId.ReplaceAllString(newContent, "")
			newContent = reResName.ReplaceAllString(newContent, "")
			newContent = rePkgName.ReplaceAllString(newContent, "")
			newContent = reKnownSub.ReplaceAllString(newContent, "")

			// Clean up double spaces if any on the go generate line
			reSpaces := regexp.MustCompile(`(go run ../../tools/generator-tests resourceidentity)[ \t]+`)
			newContent = reSpaces.ReplaceAllString(newContent, "$1 ")
			newContent = strings.ReplaceAll(newContent, "  ", " ")

			// Inject -parent-id if we inferred a virtual identity
			if d.ParentID != "" {
				// strip out the old compare-values first just to be sure
				reComp := regexp.MustCompile(`[ \t]*-compare-values[ \t]+"[^"]+"`)
				newContent = reComp.ReplaceAllString(newContent, "")

				// Check if -parent-id is already there
				if !strings.Contains(newContent, "-parent-id") {
					newContent = strings.ReplaceAll(newContent, "generator-tests resourceidentity", fmt.Sprintf("generator-tests resourceidentity -parent-id \"%s\"", d.ParentID))
				}
			}

			// Clean up any trailing space before a newline on the go generate line
			reTrailing := regexp.MustCompile(`(go run ../../tools/generator-tests resourceidentity.*?)[ \t]+\n`)
			newContent = reTrailing.ReplaceAllString(newContent, "$1\n")

			if err := os.WriteFile(goFile, []byte(newContent), 0o644); err != nil {
				return fmt.Errorf("failed writing %s for self-healing: %w", goFile, err)
			}

			// Run gofumpt
			cmd := exec.CommandContext(context.Background(), "gofumpt", "-w", goFile)
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed running gofumpt on %s: %w", goFile, err)
			}
		}
	}

	return nil
}
