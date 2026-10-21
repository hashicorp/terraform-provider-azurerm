// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_admin_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/adminrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.AdminRules

			id, err := adminrules.ParseSecurityAdminConfigurationRuleCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerAdminRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
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

			if metadata.ResourceData.HasChange("action") {
				properties.Access = model.Action
			}

			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Description = &model.Description
				} else {
					properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("destination_port_ranges") {
				properties.DestinationPortRanges = &model.DestinationPortRanges
			}

			if metadata.ResourceData.HasChange("destination") {
				properties.Destinations = expandAddressPrefixItemModel(model.Destinations)
			}

			if metadata.ResourceData.HasChange("direction") {
				properties.Direction = model.Direction
			}

			if metadata.ResourceData.HasChange("priority") {
				properties.Priority = model.Priority
			}

			if metadata.ResourceData.HasChange("protocol") {
				properties.Protocol = model.Protocol
			}

			if metadata.ResourceData.HasChange("source_port_ranges") {
				properties.SourcePortRanges = &model.SourcePortRanges
			}

			if metadata.ResourceData.HasChange("source") {
				properties.Sources = expandAddressPrefixItemModel(model.Sources)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, rule); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
