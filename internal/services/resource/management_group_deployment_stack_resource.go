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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatmanagementgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.ResourceWithUpdate = ManagementGroupDeploymentStackResource{}

type ManagementGroupDeploymentStackResource struct{}

type ManagementGroupDeploymentStackModel struct {
	Name                  string                  `tfschema:"name"`
	ManagementGroupId     string                  `tfschema:"management_group_id"`
	Location              string                  `tfschema:"location"`
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

func (r ManagementGroupDeploymentStackResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

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

		"parameters_content": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			StateFunc:        utils.NormalizeJson,
			ValidateFunc:     validation.StringIsJSON,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 256),
		},

		"action_on_unmanage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resources": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDelete),
							string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDetach),
						}, false),
					},

					"resource_groups": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDetach),
						ValidateFunc: validation.StringInSlice([]string{
							string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDelete),
							string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDetach),
						}, false),
					},

					"management_groups": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDetach),
						ValidateFunc: validation.StringInSlice([]string{
							string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDelete),
							string(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnumDetach),
						}, false),
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(deploymentstacksatmanagementgroup.DenySettingsModeNone),
							string(deploymentstacksatmanagementgroup.DenySettingsModeDenyDelete),
							string(deploymentstacksatmanagementgroup.DenySettingsModeDenyWriteAndDelete),
						}, false),
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

		"tags": commonschema.Tags(),
	}
}

func (r ManagementGroupDeploymentStackResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"output_content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"deployment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"duration": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagementGroupDeploymentStackResource) ModelObject() interface{} {
	return &ManagementGroupDeploymentStackModel{}
}

func (r ManagementGroupDeploymentStackResource) ResourceType() string {
	return "azurerm_management_group_deployment_stack"
}

func (r ManagementGroupDeploymentStackResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksManagementGroupClient

			var model ManagementGroupDeploymentStackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := deploymentstacksatmanagementgroup.NewProviders2DeploymentStackID(model.ManagementGroupId, model.Name)

			existing, err := client.DeploymentStacksGetAtManagementGroup(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := deploymentstacksatmanagementgroup.DeploymentStackProperties{
				ActionOnUnmanage: expandManagementGroupActionOnUnmanage(model.ActionOnUnmanage),
				DenySettings:     expandManagementGroupDenySettings(model.DenySettings),
			}

			if model.Description != "" {
				properties.Description = pointer.To(model.Description)
			}

			if model.TemplateContent != "" {
				template, err := expandTemplateDeploymentBody(model.TemplateContent)
				if err != nil {
					return fmt.Errorf("expanding `template_content`: %+v", err)
				}
				properties.Template = template
			}

			if model.TemplateSpecVersionId != "" {
				properties.TemplateLink = &deploymentstacksatmanagementgroup.DeploymentStacksTemplateLink{
					Id: pointer.To(model.TemplateSpecVersionId),
				}
			}

			if model.ParametersContent != "" {
				params, err := expandTemplateDeploymentBody(model.ParametersContent)
				if err != nil {
					return fmt.Errorf("expanding `parameters_content`: %+v", err)
				}
				deploymentParams := make(map[string]deploymentstacksatmanagementgroup.DeploymentParameter)
				for k, v := range *params {
					// ARM parameter files have format: {"paramName": {"value": "actualValue"}}
					// Extract the "value" field if it exists
					paramValue := v
					if paramMap, ok := v.(map[string]interface{}); ok {
						if val, exists := paramMap["value"]; exists {
							paramValue = val
						}
					}
					deploymentParams[k] = deploymentstacksatmanagementgroup.DeploymentParameter{
						Value: pointer.To(paramValue),
					}
				}
				properties.Parameters = pointer.To(deploymentParams)
			}

			tags := model.Tags
			if tags == nil {
				tags = make(map[string]string)
			}

			payload := deploymentstacksatmanagementgroup.DeploymentStack{
				Location:   pointer.To(location.Normalize(model.Location)),
				Properties: &properties,
				Tags:       pointer.To(tags),
			}

			if err := client.DeploymentStacksCreateOrUpdateAtManagementGroupThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagementGroupDeploymentStackResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksManagementGroupClient

			id, err := deploymentstacksatmanagementgroup.ParseProviders2DeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagementGroupDeploymentStackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := deploymentstacksatmanagementgroup.DeploymentStackProperties{
				ActionOnUnmanage: expandManagementGroupActionOnUnmanage(model.ActionOnUnmanage),
				DenySettings:     expandManagementGroupDenySettings(model.DenySettings),
			}

			if model.Description != "" {
				properties.Description = pointer.To(model.Description)
			}

			if model.TemplateContent != "" {
				template, err := expandTemplateDeploymentBody(model.TemplateContent)
				if err != nil {
					return fmt.Errorf("expanding `template_content`: %+v", err)
				}
				properties.Template = template
			}

			if model.TemplateSpecVersionId != "" {
				properties.TemplateLink = &deploymentstacksatmanagementgroup.DeploymentStacksTemplateLink{
					Id: pointer.To(model.TemplateSpecVersionId),
				}
			}

			if model.ParametersContent != "" {
				params, err := expandTemplateDeploymentBody(model.ParametersContent)
				if err != nil {
					return fmt.Errorf("expanding `parameters_content`: %+v", err)
				}
				deploymentParams := make(map[string]deploymentstacksatmanagementgroup.DeploymentParameter)
				for k, v := range *params {
					// ARM parameter files have format: {"paramName": {"value": "actualValue"}}
					// Extract the "value" field if it exists
					paramValue := v
					if paramMap, ok := v.(map[string]interface{}); ok {
						if val, exists := paramMap["value"]; exists {
							paramValue = val
						}
					}
					deploymentParams[k] = deploymentstacksatmanagementgroup.DeploymentParameter{
						Value: pointer.To(paramValue),
					}
				}
				properties.Parameters = pointer.To(deploymentParams)
			}

			tags := model.Tags
			if tags == nil {
				tags = make(map[string]string)
			}

			deploymentStack := deploymentstacksatmanagementgroup.DeploymentStack{
				Location:   pointer.To(location.Normalize(model.Location)),
				Properties: &properties,
				Tags:       pointer.To(tags),
			}

			if err := client.DeploymentStacksCreateOrUpdateAtManagementGroupThenPoll(ctx, *id, deploymentStack); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ManagementGroupDeploymentStackResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksManagementGroupClient

			id, err := deploymentstacksatmanagementgroup.ParseProviders2DeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeploymentStacksGetAtManagementGroup(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ManagementGroupDeploymentStackModel{
				Name:              id.DeploymentStackName,
				ManagementGroupId: id.ManagementGroupId,
			}

			if model := resp.Model; model != nil {
				if model.Location != nil {
					state.Location = location.Normalize(*model.Location)
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					state.ActionOnUnmanage = flattenManagementGroupActionOnUnmanage(props.ActionOnUnmanage)
					state.DenySettings = flattenManagementGroupDenySettings(props.DenySettings)

					if props.Description != nil {
						state.Description = *props.Description
					}

					// Handle template fields
					if props.TemplateLink != nil && props.TemplateLink.Id != nil {
						state.TemplateSpecVersionId = *props.TemplateLink.Id
					} else {
						// API doesn't return template in GET responses, preserve from current state
						state.TemplateContent = metadata.ResourceData.Get("template_content").(string)
					}

					// For parameters, preserve the ARM parameter wrapper format
					if props.Parameters != nil && len(*props.Parameters) > 0 {
						// Preserve the ARM parameter format: {"paramName": {"value": "..."}}
						params := make(map[string]interface{})
						for k, v := range *props.Parameters {
							if v.Value != nil {
								// Keep the wrapper format to match what's in config
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
					} else {
						// If API returns empty parameters but config has parameters_content, preserve it
						configParamsContent := metadata.ResourceData.Get("parameters_content").(string)
						if configParamsContent != "" {
							state.ParametersContent = configParamsContent
						}
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

func (r ManagementGroupDeploymentStackResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksManagementGroupClient

			id, err := deploymentstacksatmanagementgroup.ParseProviders2DeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagementGroupDeploymentStackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			options := deploymentstacksatmanagementgroup.DeploymentStacksDeleteAtManagementGroupOperationOptions{
				UnmanageActionResources: pointer.To(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].Resources)),
			}

			if model.ActionOnUnmanage[0].ResourceGroups != "" {
				options.UnmanageActionResourceGroups = pointer.To(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].ResourceGroups))
			}

			if model.ActionOnUnmanage[0].ManagementGroups != "" {
				options.UnmanageActionManagementGroups = pointer.To(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].ManagementGroups))
			}

			if err := client.DeploymentStacksDeleteAtManagementGroupThenPoll(ctx, *id, options); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ManagementGroupDeploymentStackResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deploymentstacksatmanagementgroup.ValidateProviders2DeploymentStackID
}

func expandManagementGroupActionOnUnmanage(input []ActionOnUnmanageModel) deploymentstacksatmanagementgroup.ActionOnUnmanage {
	if len(input) == 0 {
		return deploymentstacksatmanagementgroup.ActionOnUnmanage{}
	}

	v := input[0]
	result := deploymentstacksatmanagementgroup.ActionOnUnmanage{
		Resources: deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnum(v.Resources),
	}

	if v.ResourceGroups != "" {
		result.ResourceGroups = pointer.To(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnum(v.ResourceGroups))
	}

	if v.ManagementGroups != "" {
		result.ManagementGroups = pointer.To(deploymentstacksatmanagementgroup.DeploymentStacksDeleteDetachEnum(v.ManagementGroups))
	}

	return result
}

func flattenManagementGroupActionOnUnmanage(input deploymentstacksatmanagementgroup.ActionOnUnmanage) []ActionOnUnmanageModel {
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

func expandManagementGroupDenySettings(input []DenySettingsModel) deploymentstacksatmanagementgroup.DenySettings {
	if len(input) == 0 {
		return deploymentstacksatmanagementgroup.DenySettings{
			Mode: deploymentstacksatmanagementgroup.DenySettingsModeNone,
		}
	}

	v := input[0]
	result := deploymentstacksatmanagementgroup.DenySettings{
		Mode: deploymentstacksatmanagementgroup.DenySettingsMode(v.Mode),
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

func flattenManagementGroupDenySettings(input deploymentstacksatmanagementgroup.DenySettings) []DenySettingsModel {
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
