// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter_access_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeteraccessrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkSecurityPerimeterAccessRulesClient

			id, err := networksecurityperimeteraccessrules.ParseAccessRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config NetworkSecurityPerimeterAccessRuleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			if metadata.ResourceData.HasChange("address_prefixes") {
				existing.Model.Properties.AddressPrefixes = pointer.To(config.AddressPrefixes)
			}
			if metadata.ResourceData.HasChange("fqdns") {
				existing.Model.Properties.FullyQualifiedDomainNames = pointer.To(config.FullyQualifiedDomainNames)
			}
			if metadata.ResourceData.HasChange("service_tags") {
				existing.Model.Properties.ServiceTags = pointer.To(config.ServiceTags)
			}
			if metadata.ResourceData.HasChange("subscription_ids") {
				existing.Model.Properties.Subscriptions = expandAccessRuleSubscriptionIDs(config.Subscriptions)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
