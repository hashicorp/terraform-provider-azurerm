// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-10-01/deploymentscripts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DeploymentScriptKind string

const (
	AzurePowerShellKind DeploymentScriptKind = "AzurePowerShell"
	AzureCliKind        DeploymentScriptKind = "AzureCli"
)

type ResourceDeploymentScriptModel struct {
	Name                   string                             `tfschema:"name"`
	ResourceGroupName      string                             `tfschema:"resource_group_name"`
	Arguments              string                             `tfschema:"command_line"`
	Version                string                             `tfschema:"version"`
	CleanupPreference      deploymentscripts.CleanupOptions   `tfschema:"cleanup_preference"`
	ContainerSettings      []ContainerConfigurationModel      `tfschema:"container"`
	EnvironmentVariables   []EnvironmentVariableModel         `tfschema:"environment_variable"`
	ForceUpdateTag         string                             `tfschema:"force_update_tag"`
	Location               string                             `tfschema:"location"`
	PrimaryScriptUri       string                             `tfschema:"primary_script_uri"`
	RetentionInterval      string                             `tfschema:"retention_interval"`
	ScriptContent          string                             `tfschema:"script_content"`
	StorageAccountSettings []StorageAccountConfigurationModel `tfschema:"storage_account"`
	SupportingScriptUris   []string                           `tfschema:"supporting_script_uris"`
	Tags                   map[string]string                  `tfschema:"tags"`
	Timeout                string                             `tfschema:"timeout"`
	Outputs                string                             `tfschema:"outputs"`
}

type ContainerConfigurationModel struct {
	ContainerGroupName string `tfschema:"container_group_name"`
}

type EnvironmentVariableModel struct {
	Name        string `tfschema:"name"`
	SecureValue string `tfschema:"secure_value"`
	Value       string `tfschema:"value"`
}

type StorageAccountConfigurationModel struct {
	StorageAccountKey  string `tfschema:"key"`
	StorageAccountName string `tfschema:"name"`
}

type ResourceDeploymentScriptPatchModel struct {
	Tags map[string]string `tfschema:"tags"`
}

func getDeploymentScriptArguments(kind DeploymentScriptKind) map[string]*pluginsdk.Schema {
	result := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9_()-.]{0,259}[a-zA-Z0-9_()-]$`),
				`The name length must be from 1 to 260 characters. The name can only contain alphanumeric, underscore, parentheses, hyphen and period, and it cannot end with a period.`,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"retention_interval": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ISO8601DurationBetween("PT1H", "P1DT2H"),
		},

		"command_line": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cleanup_preference": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(deploymentscripts.CleanupOptionsOnSuccess),
				string(deploymentscripts.CleanupOptionsOnExpiration),
				string(deploymentscripts.CleanupOptionsAlways),
			}, false),
			Default: string(deploymentscripts.CleanupOptionsAlways),
		},

		"container": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_group_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"environment_variable": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secure_value": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Sensitive:    true,
					},

					"value": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			Set: hashEnvironmentVariables,
		},

		"force_update_tag": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.UserAssignedIdentityOptionalForceNew(),

		"primary_script_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"primary_script_uri", "script_content"},
		},

		"script_content": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"primary_script_uri", "script_content"},
		},

		"storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Sensitive:    true,
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"supporting_script_uris": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"timeout": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.ISO8601DurationBetween("PT1S", "P1D"),
			Default:      "P1D",
		},

		"tags": commonschema.Tags(),
	}

	if kind == AzurePowerShellKind {
		result["version"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceDeploymentScriptAzurePowerShellVersion,
		}
	} else {
		result["version"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceDeploymentScriptAzureCliVersion,
		}
	}

	return result
}

func getDeploymentScriptAttributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"outputs": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func updateDeploymentScript() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentScriptsClient

			id, err := deploymentscripts.ParseDeploymentScriptID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ResourceDeploymentScriptPatchModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &deploymentscripts.DeploymentScriptUpdateParameter{}

			if metadata.ResourceData.HasChange("tags") {
				tagValue := make(map[string]string)
				if model.Tags != nil {
					tagValue = model.Tags
				}

				properties.Tags = &tagValue
			}

			if _, err := client.Update(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func deleteDeploymentScript() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentScriptsClient

			id, err := deploymentscripts.ParseDeploymentScriptID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandContainerConfigurationModel(inputList []ContainerConfigurationModel) *deploymentscripts.ContainerConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := deploymentscripts.ContainerConfiguration{}

	if input.ContainerGroupName != "" {
		output.ContainerGroupName = &input.ContainerGroupName
	}

	return &output
}

func expandEnvironmentVariableModelArray(inputList []EnvironmentVariableModel) *[]deploymentscripts.EnvironmentVariable {
	var outputList []deploymentscripts.EnvironmentVariable
	for _, v := range inputList {
		input := v
		output := deploymentscripts.EnvironmentVariable{
			Name: input.Name,
		}

		if input.SecureValue != "" {
			output.SecureValue = &input.SecureValue
		}

		if input.Value != "" {
			output.Value = &input.Value
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandStorageAccountConfigurationModel(inputList []StorageAccountConfigurationModel) *deploymentscripts.StorageAccountConfiguration {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := deploymentscripts.StorageAccountConfiguration{}

	if input.StorageAccountKey != "" {
		output.StorageAccountKey = &input.StorageAccountKey
	}

	if input.StorageAccountName != "" {
		output.StorageAccountName = &input.StorageAccountName
	}

	return &output
}

func flattenContainerConfigurationModel(input *deploymentscripts.ContainerConfiguration) []ContainerConfigurationModel {
	var outputList []ContainerConfigurationModel
	if input == nil {
		return outputList
	}

	if input.ContainerGroupName != nil {
		output := ContainerConfigurationModel{
			ContainerGroupName: *input.ContainerGroupName,
		}
		outputList = append(outputList, output)
	}

	return outputList
}

func flattenEnvironmentVariableModelArray(inputList *[]deploymentscripts.EnvironmentVariable, originalList []EnvironmentVariableModel) []EnvironmentVariableModel {
	var outputList []EnvironmentVariableModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := EnvironmentVariableModel{
			Name: input.Name,
		}

		if input.Value != nil {
			output.Value = *input.Value
		}

		outputList = append(outputList, output)
	}

	originalSecureValues := make(map[string]string)
	for _, v := range originalList {
		if v.SecureValue != "" {
			originalSecureValues[v.Name] = v.SecureValue
		}
	}

	for i := range outputList {
		if v, exists := originalSecureValues[outputList[i].Name]; exists {
			outputList[i].SecureValue = v
		}
	}

	return outputList
}

func flattenStorageAccountConfigurationModel(input *deploymentscripts.StorageAccountConfiguration, originalList []StorageAccountConfigurationModel) []StorageAccountConfigurationModel {
	var outputList []StorageAccountConfigurationModel
	if input == nil {
		return outputList
	}

	output := StorageAccountConfigurationModel{}

	if input.StorageAccountName != nil {
		output.StorageAccountName = *input.StorageAccountName
	}

	if len(originalList) != 0 {
		output.StorageAccountKey = originalList[0].StorageAccountKey
	}

	return append(outputList, output)
}

func hashEnvironmentVariables(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["name"].(string))))
	return pluginsdk.HashString(buf.String())
}
