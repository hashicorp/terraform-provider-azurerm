// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
)

// details for a keyvault or a managedHSM
type vaultDetails struct {
	dataPlaneBaseUri string
	vaultID          string
	resourceGroup    string
}

type vaultCache struct {
	keyVaultsCache map[string]vaultDetails
	lock           *sync.RWMutex
}

func newVaultCache() *vaultCache {
	return &vaultCache{
		keyVaultsCache: map[string]vaultDetails{},
		lock:           &sync.RWMutex{},
	}
}

func keyvaultCacheKey(id commonids.KeyVaultId) string {
	return strings.ToLower(id.VaultName)
}

func (v *vaultCache) addKeyvaultToCache(key, id, resourceGroup, dataPlaneUri string) {
	item := vaultDetails{
		dataPlaneBaseUri: dataPlaneUri,
		vaultID:          id,
		resourceGroup:    resourceGroup,
	}

	v.lock.Lock()
	v.keyVaultsCache[key] = item
	v.lock.Unlock()
}

func (v *vaultCache) purge(key string) {
	v.lock.Lock()
	delete(v.keyVaultsCache, key)
	v.lock.Unlock()
}

// id can be either a ID of keyvault or managedHSM or a string of cached key
func (v *vaultCache) getCachedItem(key string) *vaultDetails {
	v.lock.RLock()
	item, ok := v.keyVaultsCache[key]
	v.lock.RUnlock()

	if ok {
		return &item
	}
	return nil
}

func (v *vaultCache) getCachedBaseUri(key string) *string {
	if item := v.getCachedItem(key); item != nil {
		return pointer.To(item.dataPlaneBaseUri)
	}
	return nil
}

func (v *vaultCache) getCachedID(key string) *string {
	if item := v.getCachedItem(key); item != nil {
		return pointer.To(item.vaultID)
	}

	return nil
}

func (c *Client) AddToCache(keyVaultId commonids.KeyVaultId, dataPlaneUri string) {
	c.keyvaultCache.addKeyvaultToCache(keyvaultCacheKey(keyVaultId), keyVaultId.ID(), keyVaultId.ResourceGroupName, dataPlaneUri)
}

func (c *Client) BaseUriForKeyVault(ctx context.Context, keyVaultId commonids.KeyVaultId) (*string, error) {
	if cached := c.keyvaultCache.getCachedBaseUri(keyvaultCacheKey(keyVaultId)); cached != nil {
		return cached, nil
	}

	resp, err := c.VaultsClient.Get(ctx, keyVaultId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s was not found", keyVaultId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", keyVaultId, err)
	}

	vaultUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.VaultUri != nil {
			vaultUri = *model.Properties.VaultUri
		}
	}
	if vaultUri == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.VaultUri` was nil", keyVaultId)
	}

	c.AddToCache(keyVaultId, vaultUri)
	return &vaultUri, nil
}

func (c *Client) Exists(ctx context.Context, keyVaultId commonids.KeyVaultId) (bool, error) {
	if c.keyvaultCache.getCachedItem(keyvaultCacheKey(keyVaultId)) != nil {
		return true, nil
	}

	resp, err := c.VaultsClient.Get(ctx, keyVaultId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving %s: %+v", keyVaultId, err)
	}

	vaultUri := ""
	if model := resp.Model; model != nil {
		if model.Properties.VaultUri != nil {
			vaultUri = *model.Properties.VaultUri
		}
	}
	if vaultUri == "" {
		return false, fmt.Errorf("retrieving %s: `properties.VaultUri` was nil", keyVaultId)
	}
	c.AddToCache(keyVaultId, vaultUri)

	return true, nil
}

func (c *Client) KeyVaultIDFromBaseUrl(ctx context.Context, subscriptionId commonids.SubscriptionId, keyVaultBaseUrl string) (*string, error) {
	keyVaultName, err := c.parseNameFromBaseUrl(keyVaultBaseUrl)
	if err != nil {
		return nil, err
	}

	if cached := c.keyvaultCache.getCachedID(strings.ToLower(*keyVaultName)); cached != nil {
		return cached, nil
	}

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
	opts := vaults.DefaultListBySubscriptionOperationOptions()
	results, err := c.VaultsClient.ListBySubscriptionComplete(ctx, subscriptionId, opts)
	if err != nil {
		return nil, fmt.Errorf("listing the Key Vaults within %s: %+v", subscriptionId, err)
	}
	for _, item := range results.Items {
		if item.Id == nil || item.Properties.VaultUri == nil {
			continue
		}

		// Populate the key vault into the cache
		keyVaultId, err := commonids.ParseKeyVaultIDInsensitively(*item.Id)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a Key Vault ID: %+v", *item.Id, err)
		}
		vaultUri := *item.Properties.VaultUri
		c.AddToCache(*keyVaultId, vaultUri)
	}

	// Now that the cache has been repopulated, check if we have the key vault or not
	if cached := c.keyvaultCache.getCachedID(*keyVaultName); cached != nil {
		return cached, nil
	}

	// We haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) Purge(keyVaultId commonids.KeyVaultId) {
	c.keyvaultCache.purge(keyvaultCacheKey(keyVaultId))
}

func (c *Client) parseNameFromBaseUrl(input string) (*string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	// https://the-keyvault.vault.azure.net
	// https://the-keyvault.vault.microsoftazure.de
	// https://the-keyvault.vault.usgovcloudapi.net
	// https://the-keyvault.vault.cloudapi.microsoft
	// https://the-keyvault.vault.azure.cn

	segments := strings.Split(uri.Host, ".")
	if len(segments) < 3 || segments[1] != "vault" {
		return nil, fmt.Errorf("expected a URI in the format `the-keyvault-name.vault.**` but got %q", uri.Host)
	}
	return &segments[0], nil
}
