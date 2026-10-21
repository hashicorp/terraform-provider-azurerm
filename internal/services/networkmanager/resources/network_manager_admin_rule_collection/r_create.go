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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerAdminRuleCollectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.AdminRuleCollections
			configurationId, err := adminrulecollections.ParseSecurityAdminConfigurationID(model.SecurityAdminConfigurationId)
			if err != nil {
				return err
			}

			id := adminrulecollections.NewSecurityAdminConfigurationRuleCollectionID(configurationId.SubscriptionId, configurationId.ResourceGroupName,
				configurationId.NetworkManagerName, configurationId.SecurityAdminConfigurationName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			adminRuleCollection := adminrulecollections.AdminRuleCollection{
				Properties: &adminrulecollections.AdminRuleCollectionPropertiesFormat{
					AppliesToGroups: expandNetworkManagerNetworkGroupIds(model.NetworkGroupIds),
				},
			}

			if model.Description != "" {
				adminRuleCollection.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, adminRuleCollection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandNetworkManagerNetworkGroupIds(inputList []string) []adminrulecollections.NetworkManagerSecurityGroupItem {
	outputList := make([]adminrulecollections.NetworkManagerSecurityGroupItem, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := adminrulecollections.NetworkManagerSecurityGroupItem{
			NetworkGroupId: input,
		}

		outputList = append(outputList, output)
	}

	return outputList
}
