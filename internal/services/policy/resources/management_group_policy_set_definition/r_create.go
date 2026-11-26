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

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetAtManagementGroup(ctx, id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := policysetdefinitions.PolicySetDefinition{
				Name: pointer.To(model.Name),
				Properties: &policysetdefinitions.PolicySetDefinitionProperties{
					Description: pointer.To(model.Description),
					DisplayName: pointer.To(model.DisplayName),
					PolicyType:  pointer.ToEnum[policysetdefinitions.PolicyType](model.PolicyType),
				},
			}
			props := parameters.Properties

			if model.Metadata != "" {
				expandedMetadata, err := pluginsdk.ExpandJsonFromString(model.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}

				var iMetadata any = expandedMetadata

				props.Metadata = &iMetadata
			}

			if model.Parameters != "" {
				expandedParameters, err := helpers.ExpandParameterDefinitionsValue(model.Parameters)
				if err != nil {
					return fmt.Errorf("expanding `parameters`: %+v", err)
				}
				props.Parameters = expandedParameters
			}

			if len(model.PolicyDefinitionReference) > 0 {
				expandedDefinitions, err := helpers.ExpandPolicyDefinitionReference(model.PolicyDefinitionReference, metadata)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if len(model.PolicyDefinitionGroup) > 0 {
				props.PolicyDefinitionGroups = helpers.ExpandPolicyDefinitionGroup(model.PolicyDefinitionGroup)
			}

			if _, err = client.CreateOrUpdateAtManagementGroup(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
