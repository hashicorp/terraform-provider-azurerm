// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

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
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AssignmentDataSource struct{}

var _ sdk.DataSource = AssignmentDataSource{}

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

func (AssignmentDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"scope_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				commonids.ValidateManagementGroupID,
				commonids.ValidateSubscriptionID,
				commonids.ValidateResourceGroupID,
				azure.ValidateResourceID,
			),
		},
	}
}

func (AssignmentDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"enforce": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"location": commonschema.LocationComputed(),

		"metadata": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"policy_definition_reference_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"parameters": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"policy_definition_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (AssignmentDataSource) ModelObject() interface{} {
	return &AssignmentDataSourceModel{}
}

func (AssignmentDataSource) ResourceType() string {
	return "azurerm_policy_assignment"
}

func (AssignmentDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.AssignmentsClient

			var plan AssignmentDataSourceModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := assignments.NewScopedPolicyAssignmentID(plan.ScopeId, plan.Name)
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
					model.Enforce = *mode == assignments.EnforcementModeDefault
				}
				model.Metadata = flattenJSON(pointer.From(props.Metadata))
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

func (m *AssignmentDataSourceModel) flattenNonComplianceMessages(input *[]assignments.NonComplianceMessage) {
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

func (m *AssignmentDataSourceModel) flattenParameter(input *map[string]assignments.ParameterValuesValue) error {
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
