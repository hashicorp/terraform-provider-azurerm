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

	// Populate the cache
	if err := c.populateCache(ctx, subscriptionId); err != nil {
		return nil, fmt.Errorf("populating the Key Vaults cache for %s: %+v", subscriptionId, err)
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
