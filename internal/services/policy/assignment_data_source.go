package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			id := parse.NewPolicyAssignmentId(plan.ScopeId, plan.Name)
			resp, err := client.Get(ctx, id.Scope, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := AssignmentDataSourceModel{
				Name:     id.Name,
				ScopeId:  id.Scope,
				Location: location.NormalizeNilable(resp.Location),
			}

			if err := model.flattenIdentity(resp.Identity); err != nil {
				return fmt.Errorf("flatten `identity`: %v", err)
			}

			if props := resp.AssignmentProperties; props != nil {
				if v := props.Description; v != nil {
					model.Description = *v
				}
				if v := props.DisplayName; v != nil {
					model.DisplayName = *v
				}
				model.Enforce = props.EnforcementMode == policy.EnforcementModeDefault
				model.Metadata = flattenJSON(props.Metadata)
				if v := props.NotScopes; v != nil {
					model.NotScopes = *v
				}
				model.flattenNonComplianceMessages(props.NonComplianceMessages)
				if err := model.flattenParameter(props.Parameters); err != nil {
					return fmt.Errorf("flatten `parameters`: %v", err)
				}
				if v := props.PolicyDefinitionID; v != nil {
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

func (m *AssignmentDataSourceModel) flattenIdentity(input *policy.Identity) error {
	if input == nil {
		return nil
	}
	config := identity.SystemOrUserAssignedMap{
		Type:        identity.Type(string(input.Type)),
		IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
	}

	if input.PrincipalID != nil {
		config.PrincipalId = *input.PrincipalID
	}
	if input.TenantID != nil {
		config.TenantId = *input.TenantID
	}
	for k, v := range input.UserAssignedIdentities {
		config.IdentityIds[k] = identity.UserAssignedIdentityDetails{
			ClientId:    v.ClientID,
			PrincipalId: v.PrincipalID,
		}
	}
	model, err := identity.FlattenSystemOrUserAssignedMapToModel(&config)
	if err != nil {
		return err
	}

	m.Identity = *model

	return nil
}

func (m *AssignmentDataSourceModel) flattenNonComplianceMessages(input *[]policy.NonComplianceMessage) {
	if input == nil {
		return
	}

	m.NonComplianceMessage = make([]NonComplianceMessage, len(*input))
	for i, v := range *input {
		content := ""
		if v.Message != nil {
			content = *v.Message
		}
		policyDefinitionReferenceId := ""
		if v.PolicyDefinitionReferenceID != nil {
			policyDefinitionReferenceId = *v.PolicyDefinitionReferenceID
		}
		m.NonComplianceMessage[i] = NonComplianceMessage{
			Content:                     content,
			PolicyDefinitionReferenceId: policyDefinitionReferenceId,
		}
	}
}

func (m *AssignmentDataSourceModel) flattenParameter(input map[string]*policy.ParameterValuesValue) error {
	if len(input) == 0 {
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
