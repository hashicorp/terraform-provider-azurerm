// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_subscription_pricing

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterSubscriptionPricingUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings.ParsePricingID(d.Id())
	if err != nil {
		return err
	}

	// Lock on subscription ID to prevent concurrent pricing updates across different resource types,
	// as the API only allows one pricing update per subscription at a time.
	locks.ByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())
	defer locks.UnlockByID(commonids.NewSubscriptionID(id.SubscriptionId).ID())

	apiResponse, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	update := pricings.Pricing{
		Properties: &pricings.PricingProperties{
			PricingTier: pricings.PricingTier(d.Get("tier").(string)),
		},
	}

	// When the state file contains an `extension` with `additional_extension_properties`
	// But the tf config does not, `d.Get("extension")` will contain a zero element.
	// Tracked by https://github.com/hashicorp/terraform-plugin-sdk/issues/1248
	realCfgExtensions := make([]any, 0)
	for _, e := range d.Get("extension").(*pluginsdk.Set).List() {
		v := e.(map[string]any)
		if v["name"] != "" {
			realCfgExtensions = append(realCfgExtensions, e)
		}
	}

	if len(realCfgExtensions) > 0 && update.Properties.PricingTier == pricings.PricingTierFree {
		return errors.New("extensions cannot be enabled when using free tier")
	}

	extensionsStatusFromBackend := make([]pricings.Extension, 0)
	currentlyFreeTier := false
	if apiResponse.Model != nil && apiResponse.Model.Properties != nil {
		if apiResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *apiResponse.Model.Properties.Extensions
		}

		currentlyFreeTier = apiResponse.Model.Properties.PricingTier == pricings.PricingTierFree
	}

	// Update from `free` tier to `Standard`, we need to update it to `standard` tier first without extensions
	// Then do an additional update for the `extensions`
	requiredAdditionalUpdate := false
	if d.HasChange("extension") && update.Properties.PricingTier == pricings.PricingTierStandard {
		update.Properties.Extensions = expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
		requiredAdditionalUpdate = currentlyFreeTier
	}

	updateResponse, err := client.Update(ctx, *id, update)
	if err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	// The extensions list from backend might vary after `tier` changed, thus we need to retrieve it again.
	if updateResponse.Model != nil && updateResponse.Model.Properties != nil {
		if updateResponse.Model.Properties.Extensions != nil {
			extensionsStatusFromBackend = *updateResponse.Model.Properties.Extensions
		}
	}

	if requiredAdditionalUpdate {
		update.Properties.Extensions = expandSecurityCenterSubscriptionPricingExtensions(realCfgExtensions, &extensionsStatusFromBackend)
		if _, err := client.Update(ctx, *id, update); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceSecurityCenterSubscriptionPricingRead(d, meta)
}
