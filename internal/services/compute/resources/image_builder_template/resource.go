// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

type Resource struct{}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualmachineimagetemplate.ValidateImageTemplateID
}

func (Resource) ResourceType() string {
	return "azurerm_image_builder_template"
}

type ImageBuilderTemplateCustomizer struct {
	Type                        string   `tfschema:"type"`
	FileDestinationPath         string   `tfschema:"file_destination_path"`
	FileSha256Checksum          string   `tfschema:"file_sha256_checksum"`
	FileSourceUri               string   `tfschema:"file_source_uri"`
	Name                        string   `tfschema:"name"`
	PowershellCommands          []string `tfschema:"powershell_commands"`
	PowershellRunAsSystem       bool     `tfschema:"powershell_run_as_system"`
	PowershellRunElevated       bool     `tfschema:"powershell_run_elevated"`
	PowershellScriptUri         string   `tfschema:"powershell_script_uri"`
	PowershellSha256Checksum    string   `tfschema:"powershell_sha256_checksum"`
	PowershellValidExitCodes    []int64  `tfschema:"powershell_valid_exit_codes"`
	ShellCommands               []string `tfschema:"shell_commands"`
	ShellScriptUri              string   `tfschema:"shell_script_uri"`
	ShellSha256Checksum         string   `tfschema:"shell_sha256_checksum"`
	WindowsRestartCheckCommand  string   `tfschema:"windows_restart_check_command"`
	WindowsRestartCommand       string   `tfschema:"windows_restart_command"`
	WindowsRestartTimeout       string   `tfschema:"windows_restart_timeout"`
	WindowsUpdateFilters        []string `tfschema:"windows_update_filters"`
	WindowsUpdateSearchCriteria string   `tfschema:"windows_update_search_criteria"`
	WindowsUpdateLimit          int64    `tfschema:"windows_update_limit"`
}

type ImageBuilderTemplateDistributions struct {
	ManagedImage []ImageBuilderTemplateDistributionsManagedImage `tfschema:"managed_image"`
	SharedImage  []ImageBuilderTemplateDistributionsSharedImage  `tfschema:"shared_image"`
	Vhd          []ImageBuilderTemplateDistributionsVhd          `tfschema:"vhd"`
}

type ImageBuilderTemplateDistributionsManagedImage struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	RunOutputName     string            `tfschema:"run_output_name"`
	Tags              map[string]string `tfschema:"tags"`
}

type ImageBuilderTemplateDistributionsSharedImage struct {
	Id                 string                                                       `tfschema:"id"`
	ReplicaRegions     []ImageBuilderTemplateDistributionsSharedImageReplicaRegions `tfschema:"replica_regions"`
	RunOutputName      string                                                       `tfschema:"run_output_name"`
	ExcludeFromLatest  bool                                                         `tfschema:"exclude_from_latest"`
	StorageAccountType string                                                       `tfschema:"storage_account_type"`
	Tags               map[string]string                                            `tfschema:"tags"`
	Versioning         []ImageBuilderTemplateDistributionsSharedImageVersioning     `tfschema:"versioning"`
}

type ImageBuilderTemplateDistributionsVhd struct {
	RunOutputName string            `tfschema:"run_output_name"`
	Tags          map[string]string `tfschema:"tags"`
}

type ImageBuilderTemplateDistributionsSharedImageReplicaRegions struct {
	Name string `tfschema:"name"`
}

type ImageBuilderTemplateDistributionsSharedImageVersioning struct {
	Scheme string `tfschema:"scheme"`
	Major  int64  `tfschema:"major"`
}

type ImageBuilderTemplateSourcePlatformImage struct {
	Publisher string                                        `tfschema:"publisher"`
	Offer     string                                        `tfschema:"offer"`
	Sku       string                                        `tfschema:"sku"`
	Version   string                                        `tfschema:"version"`
	Plan      []ImageBuilderTemplateSourcePlatformImagePlan `tfschema:"plan"`
}

type ImageBuilderTemplateSourcePlatformImagePlan struct {
	Name      string `tfschema:"name"`
	Product   string `tfschema:"product"`
	Publisher string `tfschema:"publisher"`
}

// This is to serve input validation for the customizer block.
var (
	fieldsOfFileCustomizer           = []string{"file_source_uri", "file_sha256_checksum", "file_destination_path"}
	fieldsOfPowerShellCustomizer     = []string{"powershell_commands", "powershell_run_as_system", "powershell_run_elevated", "powershell_script_uri", "powershell_sha256_checksum", "powershell_valid_exit_codes"}
	fieldsOfShellCustomizer          = []string{"shell_commands", "shell_script_uri", "shell_sha256_checksum"}
	fieldsOfWindowsRestartCustomizer = []string{"windows_restart_check_command", "windows_restart_command", "windows_restart_timeout"}
	fieldsOfWindowsUpdateCustomizer  = []string{"windows_update_filters", "windows_update_search_criteria", "windows_update_limit"}
)

func distributionRunOutputNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringMatch(
			regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
			"Run output name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
		),
	}
}

// In the `customizer` block, when a certain `type` say "File" is specified, do not allow users to specify fields that do not belong to `File`.
// Because once users specify those irrelevant fields, those values won't be sent to the backend service and they won't return from GET,
// thus next time there will be diff shown which will force creating a new resource. So at the very beginning forbid users from doing this.
func validateImageTemplateCustomizerInputForFileType(d *schema.ResourceData, customizer ImageBuilderTemplateCustomizer, index int) error {
	sourceUri := customizer.FileSourceUri
	if sourceUri == "" {
		return fmt.Errorf("`file_source_uri` must be specified if the customizer type is File")
	}

	destinationPath := customizer.FileDestinationPath
	if destinationPath == "" {
		return fmt.Errorf("`destination_path` must be specified if the customizer type is File")
	}

	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "File"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForPowerShellType(d *schema.ResourceData, customizer ImageBuilderTemplateCustomizer, index int) error {
	scriptUri := customizer.PowershellScriptUri
	commandsRaw := customizer.PowershellCommands

	if (scriptUri == "" && len(commandsRaw) == 0) ||
		(scriptUri != "" && len(commandsRaw) > 0) {
		return fmt.Errorf("exactly one of `powershell_script_uri` and `powershell_commands` must be specified if the customizer type is PowerShell")
	}

	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "PowerShell"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForShellType(d *schema.ResourceData, customizer ImageBuilderTemplateCustomizer, index int) error {
	scriptUri := customizer.ShellScriptUri
	commandsRaw := customizer.ShellCommands

	if (scriptUri == "" && len(commandsRaw) == 0) ||
		(scriptUri != "" && len(commandsRaw) > 0) {
		return fmt.Errorf("exactly one of `shell_script_uri` and `shell_commands` must be specified if the customizer type is Shell")
	}

	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "Shell"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForWindowsRestartType(d *schema.ResourceData, index int) error {
	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "WindowsRestart"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForWindowsUpdateType(d *schema.ResourceData, index int) error {
	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "WindowsUpdate"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizer(d *schema.ResourceData, exclude []string, index int, currentType string) error {
	if len(exclude) == 0 {
		return nil
	}

	for _, excludeField := range exclude {
		if _, ok := d.GetOk(fmt.Sprintf("customizer.%d.%s", index, excludeField)); ok {
			return fmt.Errorf("`%s` should not be set in the `customizer` block when the customizer type is `%s`", excludeField, currentType)
		}
	}

	return nil
}

func imageBuilderTemplateManagedImageNameAndResourceGroupName(input string) (string, string, error) {
	var imageName string
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return "", "", err
	}

	if imageName, err = id.PopSegment("images"); err != nil {
		return "", "", err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return "", "", err
	}

	return imageName, id.ResourceGroup, nil
}
