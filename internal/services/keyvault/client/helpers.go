// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/vaults"
)

var (
	keyVaultsCache = map[string]keyVaultDetails{}
	keysmith       = &sync.RWMutex{}
	lock           = map[string]*sync.RWMutex{}
)

type keyVaultDetails struct {
	keyVaultId       string
	dataPlaneBaseUri string
	resourceGroup    string
}

func (c *Client) AddToCache(keyVaultId commonids.KeyVaultId, dataPlaneUri string) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.VaultName)
	keysmith.Lock()
	keyVaultsCache[cacheKey] = keyVaultDetails{
		keyVaultId:       keyVaultId.ID(),
		dataPlaneBaseUri: dataPlaneUri,
		resourceGroup:    keyVaultId.ResourceGroupName,
	}
	keysmith.Unlock()
}

func (c *Client) BaseUriForKeyVault(ctx context.Context, keyVaultId commonids.KeyVaultId) (*string, error) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.VaultName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if v, ok := keyVaultsCache[cacheKey]; ok {
		return &v.dataPlaneBaseUri, nil
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
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.VaultName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if _, ok := keyVaultsCache[cacheKey]; ok {
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

	cacheKey := c.cacheKeyForKeyVault(*keyVaultName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	// Check the cache to determine if we have an entry for this key vault
	if v, ok := keyVaultsCache[cacheKey]; ok {
		return &v.keyVaultId, nil
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
	if v, ok := keyVaultsCache[cacheKey]; ok {
		return &v.keyVaultId, nil
	}

	// We haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) Purge(keyVaultId commonids.KeyVaultId) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.VaultName)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	delete(keyVaultsCache, cacheKey)
	lock[cacheKey].Unlock()
}

func (c *Client) cacheKeyForKeyVault(name string) string {
	return strings.ToLower(name)
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
