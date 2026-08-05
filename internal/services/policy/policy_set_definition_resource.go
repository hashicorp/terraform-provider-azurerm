// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PolicySetDefinitionResource struct{}

type PolicySetDefinitionResourceModel struct {
	Name                      string                           `tfschema:"name"`
	PolicyType                string                           `tfschema:"policy_type"`
	DisplayName               string                           `tfschema:"display_name"`
	Description               string                           `tfschema:"description"`
	Metadata                  string                           `tfschema:"metadata"`
	Parameters                string                           `tfschema:"parameters"`
	PolicyDefinitionReference []PolicyDefinitionReferenceModel `tfschema:"policy_definition_reference"`
	PolicyDefinitionGroup     []PolicyDefinitionGroupModel     `tfschema:"policy_definition_group"`
}

var (
	_ sdk.ResourceWithUpdate         = PolicySetDefinitionResource{}
	_ sdk.ResourceWithStateMigration = PolicySetDefinitionResource{}
	_ sdk.ResourceWithCustomizeDiff  = PolicySetDefinitionResource{}
)

func (r PolicySetDefinitionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.PolicySetDefinitionV0ToV1{},
		},
	}
}

func (r PolicySetDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"policy_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(policysetdefinitions.PossibleValuesForPolicyType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"metadata": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: policySetDefinitionsMetadataDiffSuppressFunc,
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"policy_definition_reference": policyDefinitionReferenceSchema(),

		"policy_definition_group": policyDefinitionGroupSchema(),
	}
}

func (r PolicySetDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PolicySetDefinitionResource) ModelObject() interface{} {
	return &PolicySetDefinitionResourceModel{}
}

func (r PolicySetDefinitionResource) ResourceType() string {
	return "azurerm_policy_set_definition"
}

func (r PolicySetDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model PolicySetDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				resp, _, err := getPolicySetDefinition(ctx, client, id)
				if err != nil && !response.WasNotFound(resp) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(resp) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := policysetdefinitions.PolicySetDefinition{
				Name: pointer.To(model.Name),
				Properties: &policysetdefinitions.PolicySetDefinitionProperties{
					Description: pointer.To(model.Description),
					DisplayName: pointer.To(model.DisplayName),
					PolicyType:  pointer.To(policysetdefinitions.PolicyType(model.PolicyType)),
				},
			}

			props := parameters.Properties
			if model.Metadata != "" {
				expandedMetadata, err := pluginsdk.ExpandJsonFromString(model.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}

				var iMetadata interface{} = expandedMetadata

				props.Metadata = &iMetadata
			}

			if model.Parameters != "" {
				expandedParameters, err := expandParameterDefinitionsValue(model.Parameters)
				if err != nil {
					return fmt.Errorf("expanding `parameters`: %+v", err)
				}
				props.Parameters = expandedParameters
			}

			if len(model.PolicyDefinitionReference) > 0 {
				expandedDefinitions, err := expandPolicyDefinitionReference(model.PolicyDefinitionReference, metadata)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if len(model.PolicyDefinitionGroup) > 0 {
				props.PolicyDefinitionGroups = expandPolicyDefinitionGroup(model.PolicyDefinitionGroup)
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PolicySetDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, model, err := getPolicySetDefinition(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PolicySetDefinitionResourceModel{
				Name: id.PolicySetDefinitionName,
			}

			if model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.DisplayName = pointer.From(props.DisplayName)
					state.PolicyType = string(pointer.From(props.PolicyType))

					if v, ok := pointer.From(props.Metadata).(map[string]interface{}); ok {
						flattenedMetadata, err := pluginsdk.FlattenJsonToString(v)
						if err != nil {
							return fmt.Errorf("flattening `metadata`: %+v", err)
						}
						state.Metadata = flattenedMetadata
					}

					flattenedParameters, err := flattenParameterDefinitionsValue(props.Parameters)
					if err != nil {
						return fmt.Errorf("flattening `parameters`: %+v", err)
					}
					state.Parameters = flattenedParameters

					flattenedDefinitions, err := flattenPolicyDefinitionReference(props.PolicyDefinitions)
					if err != nil {
						return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
					}
					state.PolicyDefinitionReference = flattenedDefinitions

					state.PolicyDefinitionGroup = flattenPolicyDefinitionGroup(props.PolicyDefinitionGroups)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PolicySetDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config PolicySetDefinitionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, model, err := getPolicySetDefinition(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}
			props := model.Properties

			if metadata.ResourceData.HasChange("display_name") {
				props.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = pointer.To(config.Description)
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
					expandedParameters, err := expandParameterDefinitionsValue(config.Parameters)
					if err != nil {
						return fmt.Errorf("expanding `parameters`: %+v", err)
					}
					props.Parameters = expandedParameters
				}
			}

			if metadata.ResourceData.HasChange("policy_definition_reference") {
				expandedDefinitions, err := expandPolicyDefinitionReference(config.PolicyDefinitionReference, metadata)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if metadata.ResourceData.HasChange("policy_definition_group") {
				props.PolicyDefinitionGroups = expandPolicyDefinitionGroup(config.PolicyDefinitionGroup)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PolicySetDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PolicySetDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return policysetdefinitions.ValidateProviderPolicySetDefinitionID
}

func getPolicySetDefinition(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, id policysetdefinitions.ProviderPolicySetDefinitionId) (*http.Response, *policysetdefinitions.PolicySetDefinition, error) {
	resp, err := client.GetBuiltIn(ctx, policysetdefinitions.NewPolicySetDefinitionID(id.PolicySetDefinitionName), policysetdefinitions.DefaultGetBuiltInOperationOptions())
	if response.WasNotFound(resp.HttpResponse) {
		resp, err := client.Get(ctx, id, policysetdefinitions.DefaultGetOperationOptions())
		return resp.HttpResponse, resp.Model, err
	}

	return resp.HttpResponse, resp.Model, err
}

func (r PolicySetDefinitionResource) CustomizeDiff() sdk.ResourceFunc {
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

					oldParameters, err := expandParameterDefinitionsValue(oldParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					newParameters, err := expandParameterDefinitionsValue(newParametersString)
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
