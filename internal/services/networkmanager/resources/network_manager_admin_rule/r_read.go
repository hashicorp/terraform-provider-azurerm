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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.AdminRules

			id, err := adminrules.ParseSecurityAdminConfigurationRuleCollectionRuleID(metadata.ResourceData.Id())
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

			var rule adminrules.AdminRule
			if adminRule, ok := existing.Model.(adminrules.AdminRule); ok {
				rule = adminRule
			}

			if rule.Properties == nil {
				return fmt.Errorf("retrieving %s: property was nil", *id)
			}

			properties := rule.Properties

			state := ManagerAdminRuleModel{
				Action: properties.Access,
				Name:   id.RuleName,
				NetworkRuleCollectionId: adminrules.NewSecurityAdminConfigurationRuleCollectionID(id.SubscriptionId, id.ResourceGroupName,
					id.NetworkManagerName, id.SecurityAdminConfigurationName, id.RuleCollectionName).ID(),
				Destinations: flattenAddressPrefixItemModel(properties.Destinations),
				Direction:    properties.Direction,
				Priority:     properties.Priority,
				Protocol:     properties.Protocol,
				Sources:      flattenAddressPrefixItemModel(properties.Sources),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.DestinationPortRanges != nil {
				state.DestinationPortRanges = *properties.DestinationPortRanges
			}

			if properties.SourcePortRanges != nil {
				state.SourcePortRanges = *properties.SourcePortRanges
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenAddressPrefixItemModel(inputList *[]adminrules.AddressPrefixItem) []AddressPrefixItemModel {
	if inputList == nil {
		return []AddressPrefixItemModel{}
	}

	outputList := make([]AddressPrefixItemModel, 0, len(*inputList))
	for _, input := range *inputList {
		output := AddressPrefixItemModel{}

		if input.AddressPrefixType != nil {
			output.AddressPrefixType = *input.AddressPrefixType
		}

		if input.AddressPrefix != nil {
			output.AddressPrefix = *input.AddressPrefix
		}

		outputList = append(outputList, output)
	}

	return outputList
}
