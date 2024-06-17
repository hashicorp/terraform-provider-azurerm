// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	vaults20230701 "github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/vaults"
	resources20151101 "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2015-11-01/resources"
)

func (c *Client) populateCache(ctx context.Context, subscriptionId commonids.SubscriptionId) error {
	// Pull out the list of Key Vaults available within the Subscription to re-populate the cache
	//
	// Whilst we've historically used the Resources API to query the single Key Vault in question
	// this endpoint has caching related issues - and whilst the ResourceGraph API has been suggested
	// as an alternative that fixes this, we've seen similar caching issues there.
	// Therefore, we're falling back on querying all the Key Vaults within the specified Subscription, which
	// comes from the `KeyVault` Resource Provider rather than the `Resources` Resource Provider - which
	// is an approach we've used previously, but now with better caching.
	//
	// Whilst querying ALL Key Vaults within a Subscription IS excessive where only a single Key Vault
	// is used - having the cache populated (one-time, per Provider launch) should alleviate problems
	// in Terraform Configurations defining a large number of Key Vault items.
	//
	// @tombuildsstuff: I vaguely recall the `ListBySubscription` API having a low rate limit (5x/second?)
	// however the rate-limits defined here seem to apply only to Managed HSMs and not Key Vaults?
	// https://learn.microsoft.com/en-us/azure/key-vault/general/service-limits
	//
	// Finally, it's worth noting that we intentionally List ALL the Key Vaults within a Subscription
	// to be able to cache ALL of them - prior to looking up the specific Key Vault we're interested
	// in from the freshly populated cache.
	// This fixes an issue in the previous implementation where the Cache was being repeatedly semi-populated
	// until the specified Key Vault was found, at which point we skipped populating the cache, which
	// affected both the `Resources` API implementation:
	// https://github.com/hashicorp/terraform-provider-azurerm/blob/3e88e5e74e12577d785f10298281b1b3c172254f/internal/services/keyvault/client/helpers.go#L133-L173
	// and the `ListBySubscription` endpoint:
	// https://github.com/hashicorp/terraform-provider-azurerm/blob/a5e728dc62e832e74d7bb0f40a79af0ae5a79e1e/azurerm/helpers/azure/key_vault.go#L42-L89
	opts := vaults20230701.DefaultListBySubscriptionOperationOptions()
	results, err := c.vaults20230701Client.ListBySubscriptionComplete(ctx, subscriptionId, opts)
	if err != nil {
		return fmt.Errorf("listing the Key Vaults within %s: %+v", subscriptionId, err)
	}
	for _, item := range results.Items {
		if item.Id == nil || item.Properties.VaultUri == nil {
			continue
		}

		// Populate the key vault into the cache
		keyVaultId, err := commonids.ParseKeyVaultIDInsensitively(*item.Id)
		if err != nil {
			return fmt.Errorf("parsing %q as a Key Vault ID: %+v", *item.Id, err)
		}
		vaultUri := *item.Properties.VaultUri
		c.AddToCache(*keyVaultId, vaultUri)
	}

	// @tombuildsstuff: now despite what I've said above, it turns out that we ALSO need to hit the Resources endpoint to populate any items that we've missed.
	// This is because the Key Vault List API appears to have a series of caching bugs where new items get added and removed from the cache, but when the cache
	// falls out of sync, it's not invalidated and reloaded - meaning that we can have totally invalid results.
	//
	// Clearly this isn't ideal, but this matches the behaviour and API version (2015-11-01) used by the Azure CLI, which should at least allow users to have
	// a consistent set of data - until the caching issues are resolved.
	resourcesOpts := resources20151101.DefaultListOperationOptions()
	resourcesOpts.Filter = pointer.To("resourceType eq 'Microsoft.KeyVault/vaults'")
	resp, err := c.resources20151101Client.ListComplete(ctx, subscriptionId, resourcesOpts)
	if err != nil {
		return fmt.Errorf("performing Resources query within %s: %+v", subscriptionId, err)
	}

	for _, item := range resp.Items {
		if item.Id == nil || item.Name == nil {
			continue
		}

		id, err := commonids.ParseKeyVaultIDInsensitively(*item.Id)
		if err != nil {
			return fmt.Errorf("parsing %q as a Key Vault ID: %+v", *item.Id, err)
		}
		cacheKey := c.cacheKeyForKeyVault(id.VaultName)
		if _, inCache := keyVaultsCache[cacheKey]; inCache {
			// don't bother caching it if we've already got it
			continue
		}

		keyVault, err := c.vaults20230701Client.Get(ctx, *id)
		if err != nil {
			if response.WasNotFound(keyVault.HttpResponse) {
				log.Printf("[DEBUG] The %s appears to be gone, this is likely an ARM Caching bug - ignoring caching it", *id)
				continue
			}
			if response.WasForbidden(keyVault.HttpResponse) {
				log.Printf("[WARN] Unable to access %s with current credentials, skipping caching it", *id)
				continue
			}
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if keyVault.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", *id)
		}
		if keyVault.Model.Properties.VaultUri == nil {
			return fmt.Errorf("retrieving %s: `model.Properties.VaultUri` was nil", *id)
		}
		dataPlaneUri := *keyVault.Model.Properties.VaultUri
		c.AddToCache(*id, dataPlaneUri)
	}

	return nil
}
