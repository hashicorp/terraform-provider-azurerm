// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_set_definition

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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

func getPolicySetDefinition(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, id policysetdefinitions.ProviderPolicySetDefinitionId) (*http.Response, *policysetdefinitions.PolicySetDefinition, error) {
	resp, err := client.GetBuiltIn(ctx, policysetdefinitions.NewPolicySetDefinitionID(id.PolicySetDefinitionName), policysetdefinitions.DefaultGetBuiltInOperationOptions())
	if response.WasNotFound(resp.HttpResponse) {
		resp, err := client.Get(ctx, id, policysetdefinitions.DefaultGetOperationOptions())
		return resp.HttpResponse, resp.Model, err
	}

	return resp.HttpResponse, resp.Model, err
}
