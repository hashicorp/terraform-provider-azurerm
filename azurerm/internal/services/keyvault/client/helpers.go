package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var keyVaultsCache map[string]keyVaultDetails
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

// nolint: deadcode unused
func (c *Client) BaseUriForKeyVault(ctx context.Context, keyVaultId parse.VaultId) (*string, error) {
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

// nolint: deadcode unused
func (c *Client) Exists(ctx context.Context, keyVaultId parse.VaultId) (bool, error) {
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

// nolint: deadcode unused
func (c *Client) KeyVaultIDFromBaseUrl(ctx context.Context, keyVaultBaseUrl string) (*string, error) {
	keyVaultName, err := c.parseNameFromBaseUrl(keyVaultBaseUrl)
	if err != nil {
		return nil, err
	}

	lock[*keyVaultName].Lock()
	defer lock[*keyVaultName].Unlock()

	list, err := c.VaultsClient.ListComplete(ctx, utils.Int32(1000))
	if err != nil {
		return nil, fmt.Errorf("failed to list Key Vaults %v", err)
	}

	// TODO: make this more efficient
	for list.NotDone() {
		v := list.Value()

		if v.ID == nil {
			return nil, fmt.Errorf("v.ID was nil")
		}

		id, err := parse.VaultID(*v.ID)
		if err != nil {
			return nil, err
		}

		// resp does not appear to contain the vault properties, so lets fetch them
		get, err := c.VaultsClient.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(get.Response) {
				if e := list.NextWithContext(ctx); e != nil {
					return nil, fmt.Errorf("failed to get next vault on KeyVault url %q : %+v", keyVaultBaseUrl, err)
				}
				continue
			}
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if get.ID == nil || get.Properties == nil || get.Properties.VaultURI == nil {
			return nil, fmt.Errorf("%s has nil ID, properties or vault URI", *id)
		}

		if keyVaultBaseUrl == *get.Properties.VaultURI {
			c.AddToCache(*id, *get.Properties.VaultURI)
			return get.ID, nil
		}

		if e := list.NextWithContext(ctx); e != nil {
			return nil, fmt.Errorf("failed to get next vault on KeyVault url %q : %+v", keyVaultBaseUrl, err)
		}
	}

	// we haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func (c *Client) Purge(keyVaultId parse.VaultId) {
	// TODO: hook this up
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
