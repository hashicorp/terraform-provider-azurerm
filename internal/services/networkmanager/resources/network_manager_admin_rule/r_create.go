// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerAdminRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.AdminRules
			ruleCollectionId, err := adminrules.ParseSecurityAdminConfigurationRuleCollectionID(model.NetworkRuleCollectionId)
			if err != nil {
				return err
			}

			id := adminrules.NewSecurityAdminConfigurationRuleCollectionRuleID(ruleCollectionId.SubscriptionId, ruleCollectionId.ResourceGroupName,
				ruleCollectionId.NetworkManagerName, ruleCollectionId.SecurityAdminConfigurationName, ruleCollectionId.RuleCollectionName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			rule := adminrules.AdminRule{
				Properties: &adminrules.AdminPropertiesFormat{
					Access:                model.Action,
					Destinations:          expandAddressPrefixItemModel(model.Destinations),
					DestinationPortRanges: &model.DestinationPortRanges,
					Direction:             model.Direction,
					Priority:              model.Priority,
					Protocol:              model.Protocol,
					SourcePortRanges:      &model.SourcePortRanges,
					Sources:               expandAddressPrefixItemModel(model.Sources),
				},
			}

			if model.Description != "" {
				rule.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, rule); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandAddressPrefixItemModel(inputList []AddressPrefixItemModel) *[]adminrules.AddressPrefixItem {
	outputList := make([]adminrules.AddressPrefixItem, 0, len(inputList))
	for _, v := range inputList {
		input := v
		output := adminrules.AddressPrefixItem{
			AddressPrefixType: &input.AddressPrefixType,
		}

		if input.AddressPrefix != "" {
			output.AddressPrefix = &input.AddressPrefix
		}

		outputList = append(outputList, output)
	}

	return &outputList
}
