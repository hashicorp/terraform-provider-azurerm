// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_group_policy_set_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

				var iMetadata any = expandedMetadata

				props.Metadata = &iMetadata
			}

			if metadata.ResourceData.HasChange("parameters") {
				props.Parameters = nil
				if config.Parameters != "" {
					expandedParameters, err := helpers.ExpandParameterDefinitionsValue(config.Parameters)
					if err != nil {
						return fmt.Errorf("expanding `parameters`: %+v", err)
					}
					props.Parameters = expandedParameters
				}
			}

			if metadata.ResourceData.HasChange("policy_definition_reference") {
				expandedDefinitions, err := helpers.ExpandPolicyDefinitionReference(config.PolicyDefinitionReference, metadata)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if metadata.ResourceData.HasChange("policy_definition_group") {
				props.PolicyDefinitionGroups = helpers.ExpandPolicyDefinitionGroup(config.PolicyDefinitionGroup)
			}

			if _, err := client.CreateOrUpdateAtManagementGroup(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
