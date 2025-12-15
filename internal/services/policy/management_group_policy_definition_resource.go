// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementGroupPolicyDefinitionResource struct{}

type ManagementGroupPolicyDefinitionResourceModel struct {
	Name              string `tfschema:"name"`
	PolicyType        string `tfschema:"policy_type"`
	ManagementGroupID string `tfschema:"management_group_id"`
	Mode              string `tfschema:"mode"`
	DisplayName       string `tfschema:"display_name"`
	Description       string `tfschema:"description"`
	Metadata          string `tfschema:"metadata"`
	Parameters        string `tfschema:"parameters"`
	PolicyRule        string `tfschema:"policy_rule"`
}

var (
	_ sdk.ResourceWithUpdate        = ManagementGroupPolicyDefinitionResource{}
	_ sdk.ResourceWithCustomizeDiff = ManagementGroupPolicyDefinitionResource{}
)

func (r ManagementGroupPolicyDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: commonids.ValidateManagementGroupID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"All",
					"Indexed",
					"Microsoft.ContainerService.Data",
					"Microsoft.CustomerLockbox.Data",
					"Microsoft.DataCatalog.Data",
					"Microsoft.KeyVault.Data",
					"Microsoft.Kubernetes.Data",
					"Microsoft.MachineLearningServices.Data",
					"Microsoft.Network.Data",
					"Microsoft.Synapse.Data",
				}, false,
			),
		},

		"policy_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(policydefinitions.PossibleValuesForPolicyType(), false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metadata": metadataSchema(),

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"policy_rule": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},
	}
}

func (r ManagementGroupPolicyDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r ManagementGroupPolicyDefinitionResource) ModelObject() interface{} {
	return &ManagementGroupPolicyDefinitionResourceModel{}
}

func (r ManagementGroupPolicyDefinitionResource) ResourceType() string {
	return "azurerm_management_group_policy_definition"
}

func (r ManagementGroupPolicyDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			var model ManagementGroupPolicyDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managementGroupID, err := commonids.ParseManagementGroupID(model.ManagementGroupID)
			if err != nil {
				return err
			}

			id := policydefinitions.NewProviders2PolicyDefinitionID(managementGroupID.GroupId, model.Name)

			existing, err := client.GetAtManagementGroup(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := policydefinitions.PolicyDefinition{
				Properties: &policydefinitions.PolicyDefinitionProperties{
					DisplayName: pointer.To(model.DisplayName),
					Description: pointer.To(model.Description),
					PolicyType:  pointer.To(policydefinitions.PolicyType(model.PolicyType)),
					Mode:        pointer.To(model.Mode),
				},
			}
			props := parameters.Properties

			if model.PolicyRule != "" {
				policyRule, err := pluginsdk.ExpandJsonFromString(model.PolicyRule)
				if err != nil {
					return fmt.Errorf("expanding `policy_rule`: %+v", err)
				}
				var iPolicyRule interface{} = policyRule
				props.PolicyRule = &iPolicyRule
			}

			if model.Metadata != "" {
				metaData, err := pluginsdk.ExpandJsonFromString(model.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}
				var iMetadata interface{} = metaData
				props.Metadata = &iMetadata
			}

			if model.Parameters != "" {
				params, err := expandParameterDefinitionsValueForPolicyDefinition(model.Parameters)
				if err != nil {
					return fmt.Errorf("expanding `parameters`: %+v", err)
				}
				props.Parameters = params
			}

			if _, err := client.CreateOrUpdateAtManagementGroup(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagementGroupPolicyDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			id, err := policydefinitions.ParseProviders2PolicyDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAtManagementGroup(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ManagementGroupPolicyDefinitionResourceModel{
				Name:              id.PolicyDefinitionName,
				ManagementGroupID: commonids.NewManagementGroupID(id.ManagementGroupName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.PolicyType = string(pointer.From(props.PolicyType))
					state.Mode = pointer.From(props.Mode)
					state.DisplayName = pointer.From(props.DisplayName)
					state.Description = pointer.From(props.Description)

					if v, ok := pointer.From(props.Metadata).(map[string]interface{}); ok {
						flattenedMetadata, err := pluginsdk.FlattenJsonToString(v)
						if err != nil {
							return fmt.Errorf("flattening `metadata`: %+v", err)
						}
						state.Metadata = flattenedMetadata
					}

					if policyRule, ok := pointer.From(props.PolicyRule).(map[string]interface{}); ok {
						flattenedPolicyRule, err := pluginsdk.FlattenJsonToString(policyRule)
						if err != nil {
							return fmt.Errorf("flattening `policy_rule`: %+v", err)
						}
						state.PolicyRule = flattenedPolicyRule

						roleIDs, _ := getPolicyRoleDefinitionIDs(flattenedPolicyRule)
						metadata.ResourceData.Set("role_definition_ids", roleIDs)
					}

					flattenedParameters, err := flattenParameterDefinitionsValueToStringForPolicyDefinition(props.Parameters)
					if err != nil {
						return fmt.Errorf("flattening `parameters`: %+v", err)
					}
					state.Parameters = flattenedParameters
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagementGroupPolicyDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			id, err := policydefinitions.ParseProviders2PolicyDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ManagementGroupPolicyDefinitionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetAtManagementGroup(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}
			props := existing.Model.Properties

			if metadata.ResourceData.HasChange("policy_type") {
				props.PolicyType = pointer.To(policydefinitions.PolicyType(config.PolicyType))
			}

			if metadata.ResourceData.HasChange("mode") {
				props.Mode = pointer.To(config.Mode)
			}

			if metadata.ResourceData.HasChange("display_name") {
				props.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("policy_rule") {
				expandedPolicyRule, err := pluginsdk.ExpandJsonFromString(config.PolicyRule)
				if err != nil {
					return fmt.Errorf("expanding `policy_rule`: %+v", err)
				}

				var iPolicyRule interface{} = expandedPolicyRule
				props.PolicyRule = &iPolicyRule
			}

			if metadata.ResourceData.HasChange("metadata") {
				expandedMetadata, err := pluginsdk.ExpandJsonFromString(config.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}

				var iMetadata interface{} = expandedMetadata

				props.Metadata = &iMetadata
			}

			if metadata.ResourceData.HasChange("parameters") {
				props.Parameters = nil
				if config.Parameters != "" {
					params, err := expandParameterDefinitionsValueForPolicyDefinition(config.Parameters)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}
					props.Parameters = params
				}
			}

			if _, err := client.CreateOrUpdateAtManagementGroup(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagementGroupPolicyDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			id, err := policydefinitions.ParseProviders2PolicyDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.DeleteAtManagementGroup(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagementGroupPolicyDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return policydefinitions.ValidateProviders2PolicyDefinitionID
}

func (r ManagementGroupPolicyDefinitionResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff.HasChange("parameters") {
				oldParametersRaw, newParametersRaw := metadata.ResourceDiff.GetChange("parameters")
				if oldParametersString := oldParametersRaw.(string); oldParametersString != "" {
					newParametersString := newParametersRaw.(string)
					if newParametersString == "" {
						return metadata.ResourceDiff.ForceNew("parameters")
					}

					oldParameters, err := expandParameterDefinitionsValueForPolicyDefinition(oldParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					newParameters, err := expandParameterDefinitionsValueForPolicyDefinition(newParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					if len(*newParameters) < len(*oldParameters) {
						return metadata.ResourceDiff.ForceNew("parameters")
					}
				}
			}

			return nil
		},
	}
}
