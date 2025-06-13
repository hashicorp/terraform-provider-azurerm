package policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagementGroupPolicySetDefinitionResource struct{}

type ManagementGroupPolicySetDefinitionResourceModel struct {
	Name                      string                           `tfschema:"name"`
	PolicyType                string                           `tfschema:"policy_type"`
	ManagementGroupID         string                           `tfschema:"management_group_id"`
	DisplayName               string                           `tfschema:"display_name"`
	Description               string                           `tfschema:"description"`
	Metadata                  string                           `tfschema:"metadata"`
	Parameters                string                           `tfschema:"parameters"`
	PolicyDefinitionReference []PolicyDefinitionReferenceModel `tfschema:"policy_definition_reference"`
	PolicyDefinitionGroup     []PolicyDefinitionGroupModel     `tfschema:"policy_definition_group"`
}

var _ sdk.ResourceWithUpdate = ManagementGroupPolicySetDefinitionResource{}

func (r ManagementGroupPolicySetDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r ManagementGroupPolicySetDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagementGroupPolicySetDefinitionResource) ModelObject() interface{} {
	return &ManagementGroupPolicySetDefinitionResourceModel{}
}

func (r ManagementGroupPolicySetDefinitionResource) ResourceType() string {
	return "azurerm_management_group_policy_set_definition"
}

func (r ManagementGroupPolicySetDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			var model ManagementGroupPolicySetDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managementGroupID, err := commonids.ParseManagementGroupID(model.ManagementGroupID)
			if err != nil {
				return err
			}

			id := policysetdefinitions.NewProviders2PolicySetDefinitionID(managementGroupID.GroupId, model.Name)

			existing, err := client.GetAtManagementGroup(ctx, id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
				expandedDefinitions, err := expandPolicyDefinitionReference(model.PolicyDefinitionReference)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if len(model.PolicyDefinitionGroup) > 0 {
				props.PolicyDefinitionGroups = expandPolicyDefinitionGroup(model.PolicyDefinitionGroup)
			}

			if _, err = client.CreateOrUpdateAtManagementGroup(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagementGroupPolicySetDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviders2PolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAtManagementGroup(ctx, *id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ManagementGroupPolicySetDefinitionResourceModel{
				Name:              id.PolicySetDefinitionName,
				ManagementGroupID: commonids.NewManagementGroupID(id.ManagementGroupName).ID(),
			}

			if model := resp.Model; model != nil {
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

func (r ManagementGroupPolicySetDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviders2PolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ManagementGroupPolicySetDefinitionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetAtManagementGroup(ctx, *id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
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
				expandedDefinitions, err := expandPolicyDefinitionReference(config.PolicyDefinitionReference)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if metadata.ResourceData.HasChange("policy_definition_group") {
				props.PolicyDefinitionGroups = expandPolicyDefinitionGroup(config.PolicyDefinitionGroup)
			}

			if _, err := client.CreateOrUpdateAtManagementGroup(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagementGroupPolicySetDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviders2PolicySetDefinitionID(metadata.ResourceData.Id())
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

func (r ManagementGroupPolicySetDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return policysetdefinitions.ValidateProviders2PolicySetDefinitionID
}
