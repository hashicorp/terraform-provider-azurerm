// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_assignment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

type AssignmentDataSourceModel struct {
	Name                 string                                     `tfschema:"name"`
	ScopeId              string                                     `tfschema:"scope_id"`
	Description          string                                     `tfschema:"description"`
	DisplayName          string                                     `tfschema:"display_name"`
	Enforce              bool                                       `tfschema:"enforce"`
	Identity             []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location             string                                     `tfschema:"location"`
	Metadata             string                                     `tfschema:"metadata"`
	NotScopes            []string                                   `tfschema:"not_scopes"`
	NonComplianceMessage []NonComplianceMessage                     `tfschema:"non_compliance_message"`
	Parameters           string                                     `tfschema:"parameters"`
	PolicyDefinitionId   string                                     `tfschema:"policy_definition_id"`
}

type NonComplianceMessage struct {
	Content                     string `tfschema:"content"`
	PolicyDefinitionReferenceId string `tfschema:"policy_definition_reference_id"`
}

func (DataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of this Policy Assignment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"scope_id": {
			Description: "The ID of the scope this Policy Assignment is assigned to. The `scope_id` can be a subscription ID, a Resource Group ID, a Management Group ID, or an ID of any resource that is assigned with a policy.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
				azure.ValidateResourceID,
			),
		},
	}
}

func (DataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "The description of this Policy Assignment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The display name of this Policy Assignment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"enforce": {
			Description: "Whether this Policy is enforced or not?",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"location": commonschema.LocationComputed(),

		"metadata": {
			Description: "A JSON mapping of any Metadata for this Policy.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"not_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"non_compliance_message": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Description: "The non-compliance message text.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"policy_definition_reference_id": {
						Description: "The ID of the Policy Definition that the non-compliance message applies to.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"parameters": {
			Description: "A JSON mapping of any Parameters for this Policy.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"policy_definition_id": {
			Description: "The ID of the assigned Policy Definition.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (DataSource) ModelObject() any {
	return &AssignmentDataSourceModel{}
}

func (DataSource) ResourceType() string {
	return "azurerm_policy_assignment"
}

func (DataSource) Description() string {
	return "Use this data source to access information about an existing Policy Assignment."
}

func (DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			var plan AssignmentDataSourceModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := policyassignments.NewScopedPolicyAssignmentID(plan.ScopeId, plan.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			respModel := resp.Model
			if respModel == nil {
				return fmt.Errorf("reading a nil model")
			}

			model := AssignmentDataSourceModel{
				Name:     id.PolicyAssignmentName,
				ScopeId:  id.Scope,
				Location: location.NormalizeNilable(respModel.Location),
			}

			if err = model.flattenIdentity(respModel.Identity); err != nil {
				return fmt.Errorf("flatten `identity`: %v", err)
			}

			if props := respModel.Properties; props != nil {
				if v := props.Description; v != nil {
					model.Description = *v
				}
				if v := props.DisplayName; v != nil {
					model.DisplayName = *v
				}
				if mode := props.EnforcementMode; mode != nil {
					model.Enforce = *mode == policyassignments.EnforcementModeDefault
				}
				model.Metadata = helpers.FlattenJSON(pointer.From(props.Metadata))
				if v := props.NotScopes; v != nil {
					model.NotScopes = *v
				}
				model.flattenNonComplianceMessages(props.NonComplianceMessages)
				if err := model.flattenParameter(props.Parameters); err != nil {
					return fmt.Errorf("flatten `parameters`: %v", err)
				}
				if v := props.PolicyDefinitionId; v != nil {
					model.PolicyDefinitionId = *v
				}
			}

			if err := metadata.Encode(&model); err != nil {
				return fmt.Errorf("encoding %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (m *AssignmentDataSourceModel) flattenNonComplianceMessages(input *[]policyassignments.NonComplianceMessage) {
	if input == nil {
		return
	}

	m.NonComplianceMessage = make([]NonComplianceMessage, len(*input))
	for i, v := range *input {
		m.NonComplianceMessage[i] = NonComplianceMessage{
			Content:                     v.Message,
			PolicyDefinitionReferenceId: pointer.From(v.PolicyDefinitionReferenceId),
		}
	}
}

func (m *AssignmentDataSourceModel) flattenParameter(input *map[string]policyassignments.ParameterValuesValue) error {
	if input == nil || len(*input) == 0 {
		return nil
	}

	result, err := json.Marshal(input)
	if err != nil {
		return err
	}

	compactJson := bytes.Buffer{}
	if err := json.Compact(&compactJson, result); err != nil {
		return err
	}

	m.Parameters = compactJson.String()
	return nil
}

func (m *AssignmentDataSourceModel) flattenIdentity(input *identity.SystemOrUserAssignedMap) error {
	model, err := identity.FlattenSystemOrUserAssignedMapToModel(input)
	if err != nil {
		return err
	}

	m.Identity = *model
	return nil
}
