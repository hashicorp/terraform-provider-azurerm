// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatsubscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubscriptionDeploymentStackResource struct{}

type SubscriptionDeploymentStackModel struct {
	Name                  string                  `tfschema:"name"`
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

	DeploymentResourceGroupName string `tfschema:"deployment_resource_group_name"`
	Location                    string `tfschema:"location"`
}

var _ sdk.ResourceWithUpdate = SubscriptionDeploymentStackResource{}

func (r SubscriptionDeploymentStackResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"action_on_unmanage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"management_groups": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnumDetach),
						ValidateFunc: validation.StringInSlice(deploymentstacksatsubscription.PossibleValuesForDeploymentStacksDeleteDetachEnum(), false),
					},

					"resource_groups": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnumDetach),
						ValidateFunc: validation.StringInSlice(deploymentstacksatsubscription.PossibleValuesForDeploymentStacksDeleteDetachEnum(), false),
					},

					"resources": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(deploymentstacksatsubscription.PossibleValuesForDeploymentStacksDeleteDetachEnum(), false),
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
						ValidateFunc: validation.StringInSlice(deploymentstacksatsubscription.PossibleValuesForDenySettingsMode(), false),
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

		"deployment_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 256),
		},

		"parameters_content": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			StateFunc:        utils.NormalizeJson,
			ValidateFunc:     validation.StringIsJSON,
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

func (r SubscriptionDeploymentStackResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r SubscriptionDeploymentStackResource) ModelObject() interface{} {
	return &SubscriptionDeploymentStackModel{}
}

func (r SubscriptionDeploymentStackResource) ResourceType() string {
	return "azurerm_subscription_deployment_stack"
}

func (r SubscriptionDeploymentStackResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksSubscriptionClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config SubscriptionDeploymentStackModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := deploymentstacksatsubscription.NewDeploymentStackID(subscriptionId, config.Name)

			existing, err := client.DeploymentStacksGetAtSubscription(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := deploymentstacksatsubscription.DeploymentStackProperties{
				ActionOnUnmanage: expandSubscriptionActionOnUnmanage(config.ActionOnUnmanage),
				DenySettings:     expandSubscriptionDenySettings(config.DenySettings),
			}

			if config.Description != "" {
				parameters.Description = pointer.To(config.Description)
			}

			if config.DeploymentResourceGroupName != "" {
				parameters.DeploymentScope = pointer.To(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, config.DeploymentResourceGroupName))
			}

			if config.TemplateContent != "" {
				template, err := expandTemplateDeploymentBody(config.TemplateContent)
				if err != nil {
					return fmt.Errorf("expanding `template_content`: %+v", err)
				}
				parameters.Template = template
			}

			if config.TemplateSpecVersionId != "" {
				parameters.TemplateLink = &deploymentstacksatsubscription.DeploymentStacksTemplateLink{
					Id: pointer.To(config.TemplateSpecVersionId),
				}
			}

			if config.ParametersContent != "" {
				params, err := expandTemplateDeploymentBody(config.ParametersContent)
				if err != nil {
					return fmt.Errorf("expanding `parameters_content`: %+v", err)
				}
				deploymentParams := make(map[string]deploymentstacksatsubscription.DeploymentParameter)
				for k, v := range *params {
					// ARM parameter files have format: {"paramName": {"value": "actualValue"}}
					// Extract the "value" field if it exists
					paramValue := v
					if paramMap, ok := v.(map[string]interface{}); ok {
						if val, exists := paramMap["value"]; exists {
							paramValue = val
						}
					}
					deploymentParams[k] = deploymentstacksatsubscription.DeploymentParameter{
						Value: pointer.To(paramValue),
					}
				}
				parameters.Parameters = pointer.To(deploymentParams)
			}

			tags := config.Tags
			if tags == nil {
				tags = make(map[string]string)
			}

			payload := deploymentstacksatsubscription.DeploymentStack{
				Location:   pointer.To(location.Normalize(config.Location)),
				Properties: &parameters,
				Tags:       pointer.To(tags),
			}

			if err := client.DeploymentStacksCreateOrUpdateAtSubscriptionThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SubscriptionDeploymentStackResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksSubscriptionClient

			id, err := deploymentstacksatsubscription.ParseDeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeploymentStacksGetAtSubscription(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := SubscriptionDeploymentStackModel{
				Name: id.DeploymentStackName,
			}

			if model := resp.Model; model != nil {
				if model.Location != nil {
					state.Location = location.Normalize(*model.Location)
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					state.ActionOnUnmanage = flattenSubscriptionActionOnUnmanage(props.ActionOnUnmanage)
					state.DenySettings = flattenSubscriptionDenySettings(props.DenySettings)

					if props.Description != nil {
						state.Description = *props.Description
					}

					if props.DeploymentScope != nil {
						scope := *props.DeploymentScope
						parts := strings.Split(scope, "/")
						if len(parts) >= 5 && parts[3] == "resourceGroups" {
							state.DeploymentResourceGroupName = parts[4]
						}
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

func (r SubscriptionDeploymentStackResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksSubscriptionClient

			id, err := deploymentstacksatsubscription.ParseDeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SubscriptionDeploymentStackModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := deploymentstacksatsubscription.DeploymentStack{
				Location:   pointer.To(location.Normalize(config.Location)),
				Properties: &deploymentstacksatsubscription.DeploymentStackProperties{},
			}

			if metadata.ResourceData.HasChange("action_on_unmanage") {
				payload.Properties.ActionOnUnmanage = expandSubscriptionActionOnUnmanage(config.ActionOnUnmanage)
			}

			if metadata.ResourceData.HasChange("deny_settings") {
				payload.Properties.DenySettings = expandSubscriptionDenySettings(config.DenySettings)
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
				payload.Properties.TemplateLink = &deploymentstacksatsubscription.DeploymentStacksTemplateLink{
					Id: pointer.To(config.TemplateSpecVersionId),
				}
			}

			if metadata.ResourceData.HasChange("parameters_content") && config.ParametersContent != "" {
				params, err := expandTemplateDeploymentBody(config.ParametersContent)
				if err != nil {
					return fmt.Errorf("expanding `parameters_content`: %+v", err)
				}
				deploymentParams := make(map[string]deploymentstacksatsubscription.DeploymentParameter)
				for k, v := range *params {
					// ARM parameter files have format: {"paramName": {"value": "actualValue"}}
					// Extract the "value" field if it exists
					paramValue := v
					if paramMap, ok := v.(map[string]interface{}); ok {
						if val, exists := paramMap["value"]; exists {
							paramValue = val
						}
					}
					deploymentParams[k] = deploymentstacksatsubscription.DeploymentParameter{
						Value: pointer.To(paramValue),
					}
				}
				payload.Properties.Parameters = pointer.To(deploymentParams)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.DeploymentStacksCreateOrUpdateAtSubscriptionThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r SubscriptionDeploymentStackResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.DeploymentStacksSubscriptionClient

			id, err := deploymentstacksatsubscription.ParseDeploymentStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SubscriptionDeploymentStackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			options := deploymentstacksatsubscription.DeploymentStacksDeleteAtSubscriptionOperationOptions{
				UnmanageActionResources: pointer.To(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].Resources)),
			}

			if model.ActionOnUnmanage[0].ResourceGroups != "" {
				options.UnmanageActionResourceGroups = pointer.To(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].ResourceGroups))
			}

			if model.ActionOnUnmanage[0].ManagementGroups != "" {
				options.UnmanageActionManagementGroups = pointer.To(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnum(model.ActionOnUnmanage[0].ManagementGroups))
			}

			if err := client.DeploymentStacksDeleteAtSubscriptionThenPoll(ctx, *id, options); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r SubscriptionDeploymentStackResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deploymentstacksatsubscription.ValidateDeploymentStackID
}

func expandSubscriptionActionOnUnmanage(input []ActionOnUnmanageModel) deploymentstacksatsubscription.ActionOnUnmanage {
	if len(input) == 0 {
		return deploymentstacksatsubscription.ActionOnUnmanage{}
	}

	v := input[0]
	result := deploymentstacksatsubscription.ActionOnUnmanage{
		Resources: deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnum(v.Resources),
	}

	if v.ResourceGroups != "" {
		result.ResourceGroups = pointer.To(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnum(v.ResourceGroups))
	}

	if v.ManagementGroups != "" {
		result.ManagementGroups = pointer.To(deploymentstacksatsubscription.DeploymentStacksDeleteDetachEnum(v.ManagementGroups))
	}

	return result
}

func flattenSubscriptionActionOnUnmanage(input deploymentstacksatsubscription.ActionOnUnmanage) []ActionOnUnmanageModel {
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

func expandSubscriptionDenySettings(input []DenySettingsModel) deploymentstacksatsubscription.DenySettings {
	if len(input) == 0 {
		return deploymentstacksatsubscription.DenySettings{
			Mode: deploymentstacksatsubscription.DenySettingsModeNone,
		}
	}

	v := input[0]
	result := deploymentstacksatsubscription.DenySettings{
		Mode: deploymentstacksatsubscription.DenySettingsMode(v.Mode),
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

func flattenSubscriptionDenySettings(input deploymentstacksatsubscription.DenySettings) []DenySettingsModel {
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
