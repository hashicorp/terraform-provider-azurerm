// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_set_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
					state.PolicyType = pointer.FromEnum(props.PolicyType)

					if v, ok := pointer.From(props.Metadata).(map[string]any); ok {
						flattenedMetadata, err := pluginsdk.FlattenJsonToString(v)
						if err != nil {
							return fmt.Errorf("flattening `metadata`: %+v", err)
						}
						state.Metadata = flattenedMetadata
					}

					flattenedParameters, err := helpers.FlattenParameterDefinitionsValue(props.Parameters)
					if err != nil {
						return fmt.Errorf("flattening `parameters`: %+v", err)
					}
					state.Parameters = flattenedParameters

					flattenedDefinitions, err := helpers.FlattenPolicyDefinitionReference(props.PolicyDefinitions)
					if err != nil {
						return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
					}
					state.PolicyDefinitionReference = flattenedDefinitions

					state.PolicyDefinitionGroup = helpers.FlattenPolicyDefinitionGroup(props.PolicyDefinitionGroups)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
