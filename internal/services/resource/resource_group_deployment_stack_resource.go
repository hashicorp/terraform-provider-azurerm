// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatresourcegroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupDeploymentStackResource struct{}

type ResourceGroupDeploymentStackModel struct {
	Name                  string                  `tfschema:"name"`
	ResourceGroupName     string                  `tfschema:"resource_group_name"`
	TemplateContent       string                  `tfschema:"template_content"`
	TemplateSpecVersionId string                  `tfschema:"template_spec_version_id"`
	ParametersContent     string                  `tfschema:"parameters_content"`
	Description           string                  `tfschema:"description"`
	ActionOnUnmanage      []ActionOnUnmanageModel `tfschema:"action_on_unmanage"`
	DenySettings          []DenySettingsModel     `tfschema:"deny_settings"`
	Tags                  map[string]string       `tfschema:"tags"`
	OutputContent         string                  `tfschema:"output_content"`
	DeploymentId          string                  `tfschema:"deployment_id"`
	Duration              string                  `tfschema:"duration"`
}

var _ sdk.ResourceWithUpdate = ResourceGroupDeploymentStackResource{}

func (r ResourceGroupDeploymentStackResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"action_on_unmanage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"management_groups": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnumDetach),
						ValidateFunc: validation.StringInSlice(deploymentstacksatresourcegroup.PossibleValuesForDeploymentStacksDeleteDetachEnum(), false),
					},

					"resource_groups": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnumDetach),
						ValidateFunc: validation.StringInSlice(deploymentstacksatresourcegroup.PossibleValuesForDeploymentStacksDeleteDetachEnum(), false),
					},

					"resources": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(deploymentstacksatresourcegroup.PossibleValuesForDeploymentStacksDeleteDetachEnum(), false),
					},
				},
			},
		},

		"deny_settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(deploymentstacksatresourcegroup.PossibleValuesForDenySettingsMode(), false),
					},

					"apply_to_child_scopes": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"excluded_actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 200,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"excluded_principals": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 5,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 256),
		},

		"parameters_content": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			StateFunc:    utils.NormalizeJson,
			ValidateFunc: validation.StringIsJSON,
		},

		"tags": commonschema.Tags(),

		"template_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ExactlyOneOf: []string{
				"template_content",
				"template_spec_version_id",
			},
			StateFunc:        utils.NormalizeJson,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"template_spec_version_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ExactlyOneOf: []string{
				"template_content",
				"template_spec_version_id",
			},
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ResourceGroupDeploymentStackResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"deployment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"duration": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"output_content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ResourceGroupDeploymentStackResource) ModelObject() interface{} {
	return &ResourceGroupDeploymentStackModel{}
}

func (r ResourceGroupDeploymentStackResource) ResourceType() string {
	return "azurerm_resource_group_deployment_stack"
}

func (r ResourceGroupDeploymentStackResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksResourceGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ResourceGroupDeploymentStackModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := deploymentstacksatresourcegroup.NewProviderDeploymentStackID(subscriptionId, config.ResourceGroupName, config.Name)

			existing, err := client.DeploymentStacksGetAtResourceGroup(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := deploymentstacksatresourcegroup.DeploymentStack{
				Properties: &deploymentstacksatresourcegroup.DeploymentStackProperties{
					ActionOnUnmanage: expandActionOnUnmanage(config.ActionOnUnmanage),
					DenySettings:     expandDenySettings(config.DenySettings),
				},
				// Tags: pointer.To(tags),
			}

			properties := parameters.Properties

			if config.Description != "" {
				properties.Description = pointer.To(config.Description)
			}

			if config.TemplateContent != "" {
				template, err := expandTemplateDeploymentBody(config.TemplateContent)
				if err != nil {
					return fmt.Errorf("expanding `template_content`: %+v", err)
				}
				properties.Template = template
			}

			if config.TemplateSpecVersionId != "" {
				properties.TemplateLink = &deploymentstacksatresourcegroup.DeploymentStacksTemplateLink{
					Id: pointer.To(config.TemplateSpecVersionId),
				}
			}

			if config.ParametersContent != "" {
				params, err := expandTemplateDeploymentBody(config.ParametersContent)
				if err != nil {
					return fmt.Errorf("expanding `parameters_content`: %+v", err)
				}
				deploymentParams := make(map[string]deploymentstacksatresourcegroup.DeploymentParameter)
				for k, v := range *params {
					// ARM parameter files have format: {"paramName": {"value": "actualValue"}}
					// Extract the "value" field if it exists
					paramValue := v
					if paramMap, ok := v.(map[string]interface{}); ok {
						if val, exists := paramMap["value"]; exists {
							paramValue = val
						}
					}
					deploymentParams[k] = deploymentstacksatresourcegroup.DeploymentParameter{
						Value: pointer.To(paramValue),
					}
				}
				properties.Parameters = pointer.To(deploymentParams)
			}

			parameters.Tags = pointer.To(config.Tags)
			if config.Tags == nil {
				parameters.Tags = pointer.To(make(map[string]string))
			}

			if err := client.DeploymentStacksCreateOrUpdateAtResourceGroupThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceGroupDeploymentStackResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksResourceGroupClient

			id, err := deploymentstacksatresourcegroup.ParseProviderDeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeploymentStacksGetAtResourceGroup(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ResourceGroupDeploymentStackModel{
				Name:              id.DeploymentStackName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					state.ActionOnUnmanage = flattenActionOnUnmanage(props.ActionOnUnmanage)
					state.DenySettings = flattenDenySettings(props.DenySettings)

					if props.Description != nil {
						state.Description = *props.Description
					}

					// Azure does not return `template` in GET responses, thus default to value in state
					state.TemplateContent = metadata.ResourceData.Get("template_content").(string)
					if props.TemplateLink != nil && props.TemplateLink.Id != nil {
						state.TemplateSpecVersionId = *props.TemplateLink.Id
					}

					configParamsContent := metadata.ResourceData.Get("parameters_content").(string)
					if configParamsContent != "" {
						state.ParametersContent = configParamsContent
					}

					// If `parameters` is empty in API, preserve the config value
					if props.Parameters != nil && len(*props.Parameters) > 0 {
						// Preserve the ARM parameter format: {"paramName": {"value": "..."}}
						params := make(map[string]interface{})
						for k, v := range *props.Parameters {
							if v.Value != nil {
								params[k] = map[string]interface{}{
									"value": *v.Value,
								}
							}
						}
						flattenedParams, err := flattenTemplateDeploymentBody(params)
						if err != nil {
							return fmt.Errorf("flattening `parameters_content`: %+v", err)
						}
						state.ParametersContent = *flattenedParams
					}

					if props.Outputs != nil {
						flattenedOutputs, err := flattenTemplateDeploymentBody(*props.Outputs)
						if err != nil {
							return fmt.Errorf("flattening `output_content`: %+v", err)
						}
						state.OutputContent = *flattenedOutputs
					}

					if props.DeploymentId != nil {
						state.DeploymentId = *props.DeploymentId
					}

					if props.Duration != nil {
						state.Duration = *props.Duration
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ResourceGroupDeploymentStackResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksResourceGroupClient

			id, err := deploymentstacksatresourcegroup.ParseProviderDeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ResourceGroupDeploymentStackModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.DeploymentStacksGetAtResourceGroup(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := pointer.From(existing.Model)

			if metadata.ResourceData.HasChange("action_on_unmanage") {
				payload.Properties.ActionOnUnmanage = expandActionOnUnmanage(config.ActionOnUnmanage)
			}

			if metadata.ResourceData.HasChange("deny_settings") {
				payload.Properties.DenySettings = expandDenySettings(config.DenySettings)
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("template_content") {
				template, err := expandTemplateDeploymentBody(config.TemplateContent)
				if err != nil {
					return fmt.Errorf("expanding `template_content`: %+v", err)
				}
				payload.Properties.Template = template
			}

			if metadata.ResourceData.HasChange("template_spec_version_id") {
				payload.Properties.TemplateLink = &deploymentstacksatresourcegroup.DeploymentStacksTemplateLink{
					Id: pointer.To(config.TemplateSpecVersionId),
				}
			}

			if metadata.ResourceData.HasChange("parameters_content") {
				params, err := expandTemplateDeploymentBody(config.ParametersContent)
				if err != nil {
					return fmt.Errorf("expanding `parameters_content`: %+v", err)
				}
				deploymentParams := make(map[string]deploymentstacksatresourcegroup.DeploymentParameter)
				for k, v := range *params {
					// ARM parameter files have format: {"paramName": {"value": "actualValue"}}
					// Extract the "value" field if it exists
					paramValue := v
					if paramMap, ok := v.(map[string]interface{}); ok {
						if val, exists := paramMap["value"]; exists {
							paramValue = val
						}
					}
					deploymentParams[k] = deploymentstacksatresourcegroup.DeploymentParameter{
						Value: pointer.To(paramValue),
					}
				}
				payload.Properties.Parameters = pointer.To(deploymentParams)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.DeploymentStacksCreateOrUpdateAtResourceGroupThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ResourceGroupDeploymentStackResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksResourceGroupClient

			id, err := deploymentstacksatresourcegroup.ParseProviderDeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ResourceGroupDeploymentStackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			options := deploymentstacksatresourcegroup.DeploymentStacksDeleteAtResourceGroupOperationOptions{
				UnmanageActionResources: pointer.To(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].Resources)),
			}

			if model.ActionOnUnmanage[0].ResourceGroups != "" {
				options.UnmanageActionResourceGroups = pointer.To(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].ResourceGroups))
			}

			if model.ActionOnUnmanage[0].ManagementGroups != "" {
				options.UnmanageActionManagementGroups = pointer.To(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].ManagementGroups))
			}

			if err := client.DeploymentStacksDeleteAtResourceGroupThenPoll(ctx, *id, options); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ResourceGroupDeploymentStackResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deploymentstacksatresourcegroup.ValidateProviderDeploymentStackID
}

func expandActionOnUnmanage(input []ActionOnUnmanageModel) deploymentstacksatresourcegroup.ActionOnUnmanage {
	if len(input) == 0 {
		return deploymentstacksatresourcegroup.ActionOnUnmanage{}
	}

	v := input[0]
	result := deploymentstacksatresourcegroup.ActionOnUnmanage{
		Resources: deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnum(v.Resources),
	}

	if v.ResourceGroups != "" {
		result.ResourceGroups = pointer.To(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnum(v.ResourceGroups))
	}

	if v.ManagementGroups != "" {
		result.ManagementGroups = pointer.To(deploymentstacksatresourcegroup.DeploymentStacksDeleteDetachEnum(v.ManagementGroups))
	}

	return result
}

func flattenActionOnUnmanage(input deploymentstacksatresourcegroup.ActionOnUnmanage) []ActionOnUnmanageModel {
	result := ActionOnUnmanageModel{
		Resources: string(input.Resources),
	}

	if input.ResourceGroups != nil {
		result.ResourceGroups = string(*input.ResourceGroups)
	}

	if input.ManagementGroups != nil {
		result.ManagementGroups = string(*input.ManagementGroups)
	}

	return []ActionOnUnmanageModel{result}
}

func expandDenySettings(input []DenySettingsModel) deploymentstacksatresourcegroup.DenySettings {
	if len(input) == 0 {
		return deploymentstacksatresourcegroup.DenySettings{
			Mode: deploymentstacksatresourcegroup.DenySettingsModeNone,
		}
	}

	v := input[0]
	result := deploymentstacksatresourcegroup.DenySettings{
		Mode: deploymentstacksatresourcegroup.DenySettingsMode(v.Mode),
	}

	if v.ApplyToChildScopes {
		result.ApplyToChildScopes = pointer.To(true)
	}

	if v.ExcludedActions != nil && len(*v.ExcludedActions) > 0 {
		result.ExcludedActions = v.ExcludedActions
	}

	if v.ExcludedPrincipals != nil && len(*v.ExcludedPrincipals) > 0 {
		result.ExcludedPrincipals = v.ExcludedPrincipals
	}

	return result
}

func flattenDenySettings(input deploymentstacksatresourcegroup.DenySettings) []DenySettingsModel {
	result := DenySettingsModel{
		Mode: string(input.Mode),
	}

	if input.ApplyToChildScopes != nil {
		result.ApplyToChildScopes = *input.ApplyToChildScopes
	}

	// Only set excluded_actions/excluded_principals if they have values
	// Don't set them at all if empty to avoid drift
	if input.ExcludedActions != nil && len(*input.ExcludedActions) > 0 {
		result.ExcludedActions = input.ExcludedActions
	}

	if input.ExcludedPrincipals != nil && len(*input.ExcludedPrincipals) > 0 {
		result.ExcludedPrincipals = input.ExcludedPrincipals
	}

	return []DenySettingsModel{result}
}
