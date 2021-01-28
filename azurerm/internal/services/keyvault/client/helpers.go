package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	resource "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/client"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var keyVaultsCache = map[string]keyVaultDetails{}
var lock = map[string]*sync.RWMutex{}

type keyVaultDetails struct {
	keyVaultId       string
	dataPlaneBaseUri string
	resourceGroup    string
}

func (c *Client) AddToCache(keyVaultId parse.VaultId, dataPlaneUri string) {
	keyVaultsCache[keyVaultId.Name] = keyVaultDetails{
		keyVaultId:       keyVaultId.ID(),
		dataPlaneBaseUri: dataPlaneUri,
		resourceGroup:    keyVaultId.ResourceGroup,
	}
}

func (c *Client) BaseUriForKeyVault(ctx context.Context, keyVaultId parse.VaultId) (*string, error) {
	if lock[keyVaultId.Name] == nil {
		lock[keyVaultId.Name] = &sync.RWMutex{}
	}
	lock[keyVaultId.Name].Lock()
	defer lock[keyVaultId.Name].Unlock()

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
	if lock[keyVaultId.Name] == nil {
		lock[keyVaultId.Name] = &sync.RWMutex{}
	}
	lock[keyVaultId.Name].Lock()
	defer lock[keyVaultId.Name].Unlock()

	if _, ok := keyVaultsCache[keyVaultId.Name]; ok {
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

func (c *Client) KeyVaultIDFromBaseUrl(ctx context.Context, resourcesClient *resource.Client, keyVaultBaseUrl string) (*string, error) {
	keyVaultName, err := c.parseNameFromBaseUrl(keyVaultBaseUrl)
	if err != nil {
		return nil, err
	}

	if lock[*keyVaultName] == nil {
		lock[*keyVaultName] = &sync.RWMutex{}
	}
	lock[*keyVaultName].Lock()
	defer lock[*keyVaultName].Unlock()

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
			if id.Name != *keyVaultName {
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
	if lock[keyVaultId.Name] == nil {
		lock[keyVaultId.Name] = &sync.RWMutex{}
	}
	lock[keyVaultId.Name].Lock()
	defer lock[keyVaultId.Name].Unlock()
	delete(keyVaultsCache, keyVaultId.Name)
}

func (c *Client) parseNameFromBaseUrl(input string) (*string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	// https://tharvey-keyvault.vault.azure.net/
	segments := strings.Split(uri.Host, ".")
	if len(segments) != 4 {
		return nil, fmt.Errorf("expected a URI in the format `vaultname.vault.azure.net` but got %q", uri.Host)
	}
	return &segments[0], nil
}
