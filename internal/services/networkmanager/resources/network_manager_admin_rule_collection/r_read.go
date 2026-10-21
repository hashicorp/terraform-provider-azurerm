// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule_collection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrulecollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.AdminRuleCollections

			id, err := adminrulecollections.ParseSecurityAdminConfigurationRuleCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerAdminRuleCollectionModel{
				Name:                         id.RuleCollectionName,
				SecurityAdminConfigurationId: adminrulecollections.NewSecurityAdminConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.SecurityAdminConfigurationName).ID(),
				NetworkGroupIds:              flattenNetworkManagerNetworkGroupIds(properties.AppliesToGroups),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenNetworkManagerNetworkGroupIds(inputList []adminrulecollections.NetworkManagerSecurityGroupItem) []string {
	outputList := make([]string, 0, len(inputList))

	for _, input := range inputList {
		outputList = append(outputList, input.NetworkGroupId)
	}

	return outputList
}
