package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2021-10-01/keyvault" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type vaultCahcher interface {
	GetVault(vault parse.Vaulter) *keyVaultDetails
	GetVaultByKey(key string) *keyVaultDetails
	AddVault(vault parse.Vaulter, dataPlaneURI string)
	Delete(vault parse.Vaulter)
}

type vaultCache struct {
	cache map[string]keyVaultDetails
	lock  sync.RWMutex
}

func (v *vaultCache) GetVault(vault parse.Vaulter) (res *keyVaultDetails) {
	return v.GetVaultByKey(vault.GetCacheKey())
}

func (v *vaultCache) GetVaultByKey(key string) (res *keyVaultDetails) {
	v.lock.RLock()
	if val, ok := v.cache[key]; ok {
		res = &val
	}
	v.lock.RUnlock()
	return
}

func (v *vaultCache) AddVault(vault parse.Vaulter, dataPlaneURI string) {
	v.lock.Lock()
	v.cache[vault.GetCacheKey()] = keyVaultDetails{
		keyVaultId:       vault.ID(),
		dataPlaneBaseUri: dataPlaneURI,
		resourceGroup:    vault.GetResourceGroup(),
	}
	v.lock.Unlock()
}

func (v *vaultCache) Delete(vault parse.Vaulter) {
	v.lock.Lock()
	delete(v.cache, vault.GetCacheKey())
	v.lock.Unlock()
}

var vaultCacheIns vaultCahcher = &vaultCache{
	cache: map[string]keyVaultDetails{},
}

type keyVaultDetails struct {
	keyVaultId       string
	dataPlaneBaseUri string
	resourceGroup    string
}

func (c *Client) AddToCache(keyVaultId parse.Vaulter, dataPlaneUri string) {
	vaultCacheIns.AddVault(keyVaultId, dataPlaneUri)
}

func (c *Client) BaseUriForKeyVault(ctx context.Context, keyVaultId parse.Vaulter) (*string, error) {
	if cached := vaultCacheIns.GetVault(keyVaultId); cached != nil {
		return &cached.dataPlaneBaseUri, nil
	}
	var (
		baseURI       *string
		subscription  = keyVaultId.GetSubscriptionID()
		resourceGroup = keyVaultId.GetResourceGroup()
		name          = keyVaultId.GetName()
		err           error
	)

	switch vault := keyVaultId.(type) {
	case parse.VaultId, *parse.VaultId:
		if subscription != c.VaultsClient.SubscriptionID {
			c.VaultsClient = c.KeyVaultClientForSubscription(subscription)
		}

		resp, err := c.VaultsClient.Get(ctx, resourceGroup, vault.GetName())
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil, fmt.Errorf("%s was not found", keyVaultId)
			}
			return nil, fmt.Errorf("retrieving %s: %+v", keyVaultId, err)
		}

		if resp.Properties == nil || resp.Properties.VaultURI == nil {
			return nil, fmt.Errorf("`properties` was nil for %s", keyVaultId)
		}

		baseURI = resp.Properties.VaultURI
	case parse.ManagedHSMId, *parse.ManagedHSMId:
		resp, err := c.ManagedHsmClient.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil, fmt.Errorf("%s was not found", keyVaultId)
			}
			return nil, fmt.Errorf("retrieving %s: %+v", keyVaultId, err)
		}

		if resp.Properties == nil || resp.Properties.HsmURI == nil {
			return nil, fmt.Errorf("`properties` was nil for %s", keyVaultId)
		}

		baseURI = resp.Properties.HsmURI
	default:
		err = fmt.Errorf("not support key vault type: %q", keyVaultId)
	}
	if baseURI != nil {
		vaultCacheIns.AddVault(keyVaultId, *baseURI)
	}
	return baseURI, err
}

func (c *Client) Exists(ctx context.Context, keyVaultId parse.Vaulter) (bool, error) {
	if cached := vaultCacheIns.GetVault(keyVaultId); cached != nil {
		return true, nil
	}

	uri, err := c.GetVaultURI(ctx, keyVaultId)
	if err != nil || uri == "" {
		return false, err
	}

	c.AddToCache(keyVaultId, uri)

	return true, nil
}

// GetVault try to get KeyVault or HSM instance of vaulter
func (c *Client) GetVault(ctx context.Context, vaulter parse.Vaulter) (
	vault *keyvault.Vault,
	hsm *keyvault.ManagedHsm,
	err error) {

	switch vaulter.Type() {
	case parse.VaultTypeDefault:
		vaultIns, err := c.VaultsClient.Get(ctx, vaulter.GetResourceGroup(), vaulter.GetName())
		return &vaultIns, nil, err
	case parse.VaultTypeMHSM:
		hsmIns, err := c.ManagedHsmClient.Get(ctx, vaulter.GetResourceGroup(), vaulter.GetName())
		return nil, &hsmIns, err
	}
	return nil, nil, fmt.Errorf("not supported type: %s", vaulter.Type())
}

// GetVaultURI Get uri of key vault or uri of hsm
func (c *Client) GetVaultURI(ctx context.Context, vaulter parse.Vaulter) (uri string, err error) {
	vault, hsm, err := c.GetVault(ctx, vaulter)
	if vault != nil && vault.Properties != nil && vault.Properties.VaultURI != nil {
		uri = *vault.Properties.VaultURI
	} else if hsm != nil && hsm.Properties != nil && hsm.Properties.HsmURI != nil {
		uri = *hsm.Properties.HsmURI
	}
	return uri, err
}

func (c *Client) KeyVaultIDFromBaseUrl(ctx context.Context,
	resourcesClient *resourcesClient.Client,
	keyVaultBaseUrl string) (
	*string, error) {

	keyVaultName, vaultType, err := c.parseNameFromBaseUrl(keyVaultBaseUrl)
	if err != nil {
		return nil, err
	}

	if cached := vaultCacheIns.GetVaultByKey(parse.MakeCacheKey(vaultType, pointer.From(keyVaultName))); cached != nil {
		return &cached.keyVaultId, nil
	}

	filter := fmt.Sprintf("resourceType eq 'Microsoft.KeyVault/vaults' and name eq '%s'", *keyVaultName)
	isMHSMVault := vaultType == parse.VaultTypeMHSM
	if isMHSMVault {
		filter = fmt.Sprintf("resourceType eq 'Microsoft.KeyVault/managedHSMs' and name eq '%s'", *keyVaultName)
	}
	result, err := resourcesClient.ResourcesClient.List(ctx, filter, "", utils.Int32(5))
	if err != nil {
		return nil, fmt.Errorf("listing resources matching %q: %+v", filter, err)
	}

	for result.NotDone() {
		for _, v := range result.Values() {
			if v.ID == nil {
				continue
			}

			id, err := parse.NewVaulterFromString(*v.ID)
			if err != nil {
				return nil, fmt.Errorf("parsing %q: %+v", *v.ID, err)
			}
			if !strings.EqualFold(id.GetName(), *keyVaultName) {
				continue
			}

			vaultURI, err := c.GetVaultURI(ctx, id)
			if err != nil {
				return nil, err
			}

			c.AddToCache(id, vaultURI)
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
	vaultCacheIns.Delete(keyVaultId)
}

func (c *Client) parseNameFromBaseUrl(input string) (*string, parse.VaultType, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, "", err
	}

	// https://the-keyvault.vault.azure.net
	// https://the-keyvault.vault.microsoftazure.de
	// https://the-keyvault.vault.usgovcloudapi.net
	// https://the-keyvault.vault.cloudapi.microsoft
	// https://the-keyvault.vault.azure.cn

	segments := strings.Split(uri.Host, ".")
	if len(segments) < 3 || !parse.IsValidValtType(segments[1]) {
		return nil, "", fmt.Errorf("expected a URI in the format `the-keyvault-name.vault.**` or `the-keyvault-name.managedhsm.**` but got %q", uri.Host)
	}
	return &segments[0], parse.VaultType(segments[1]), nil
}
