package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	resourcesClient "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var keyVaultsCache = map[string]keyVaultDetails{}
var keysmith = &sync.RWMutex{}
var lock = map[string]*sync.RWMutex{}

type keyVaultDetails struct {
	keyVaultId       string
	dataPlaneBaseUri string
	resourceGroup    string
}

func (c *Client) AddToCache(keyVaultId parse.VaultId, dataPlaneUri string) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.Name)
	keysmith.Lock()
	keyVaultsCache[cacheKey] = keyVaultDetails{
		keyVaultId:       keyVaultId.ID(),
		dataPlaneBaseUri: dataPlaneUri,
		resourceGroup:    keyVaultId.ResourceGroup,
	}
	keysmith.Unlock()
}

func (c *Client) BaseUriForKeyVault(ctx context.Context, keyVaultId parse.VaultId) (*string, error) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.Name)
	keysmith.Lock()
	if lock[cacheKey] == nil {
		lock[cacheKey] = &sync.RWMutex{}
	}
	keysmith.Unlock()
	lock[cacheKey].Lock()
	defer lock[cacheKey].Unlock()

	if keyVaultId.SubscriptionId != c.VaultsClient.SubscriptionID {
		c.VaultsClient = c.KeyVaultClientForSubscription(keyVaultId.SubscriptionId)
	}

	resp, err := c.VaultsClient.Get(ctx, keyVaultId.ResourceGroup, keyVaultId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, fmt.Errorf("%s was not found", keyVaultId)
		}
		return nil, fmt.Errorf("retrieving %s: %+v", keyVaultId, err)
	}

	if resp.Properties == nil || resp.Properties.VaultURI == nil {
		return nil, fmt.Errorf("`properties` was nil for %s", keyVaultId)
	}

	return resp.Properties.VaultURI, nil
}

func (c *Client) Exists(ctx context.Context, keyVaultId parse.VaultId) (bool, error) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.Name)
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

	resp, err := c.VaultsClient.Get(ctx, keyVaultId.ResourceGroup, keyVaultId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving %s: %+v", keyVaultId, err)
	}

	if resp.Properties == nil || resp.Properties.VaultURI == nil {
		return false, fmt.Errorf("`properties` was nil for %s", keyVaultId)
	}

	c.AddToCache(keyVaultId, *resp.Properties.VaultURI)

	return true, nil
}

func (c *Client) KeyVaultIDFromBaseUrl(ctx context.Context, resourcesClient *resourcesClient.Client, keyVaultBaseUrl string) (*string, error) {
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

	if v, ok := keyVaultsCache[cacheKey]; ok {
		return &v.keyVaultId, nil
	}

	filter := fmt.Sprintf("resourceType eq 'Microsoft.KeyVault/vaults' and name eq '%s'", *keyVaultName)
	result, err := resourcesClient.ResourcesClient.List(ctx, filter, "", utils.Int32(5))
	if err != nil {
		return nil, fmt.Errorf("listing resources matching %q: %+v", filter, err)
	}

	for result.NotDone() {
		for _, v := range result.Values() {
			if v.ID == nil {
				continue
			}

			id, err := parse.VaultID(*v.ID)
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %+v", *v.ID, err)
			}
			if !strings.EqualFold(id.Name, *keyVaultName) {
				continue
			}

			props, err := c.VaultsClient.Get(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if props.Properties == nil || props.Properties.VaultURI == nil {
				return nil, fmt.Errorf("retrieving %s: `properties.VaultUri` was nil", *id)
			}

			c.AddToCache(*id, *props.Properties.VaultURI)
			return utils.String(id.ID()), nil
		}

		if err := result.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("iterating over results: %+v", err)
		}
	}

	// we haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) Purge(keyVaultId parse.VaultId) {
	cacheKey := c.cacheKeyForKeyVault(keyVaultId.Name)
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
