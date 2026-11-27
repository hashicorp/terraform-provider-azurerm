// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_subscription_pricing

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterSubscriptionPricingDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings.ParsePricingID(d.Id())
	if err != nil {
		return err
	}

	// Lock on subscription ID to prevent concurrent pricing updates across different resource types,
	// as the API only allows one pricing update per subscription at a time.
	locks.ByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())
	defer locks.UnlockByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())

	pricing := pricings.Pricing{
		Properties: &pricings.PricingProperties{
			PricingTier: pricings.PricingTierFree,
		},
	}

	if _, err := client.Update(ctx, *id, pricing); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil
}
