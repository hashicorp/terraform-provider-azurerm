package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var (
	keyVaultsCache = map[string]*keyVaultDetails{}
	KeyVaultLock   = sync.RWMutex{}
)

type keyVaultDetails struct {
	ID            string
	ResourceGroup string
	Name          string
	Properties    *keyvault.VaultProperties
	Tags          map[string]*string
}

func (client Client) RemoveKeyVaultFromCache(keyVaultUrl string) {
	KeyVaultLock.Lock()
	delete(keyVaultsCache, keyVaultUrl)
	KeyVaultLock.Unlock()
}

func (client Client) FindKeyVault(ctx context.Context, keyVaultUrl string) (*keyVaultDetails, error) {
	KeyVaultLock.Lock()
	defer KeyVaultLock.Unlock()

	if existing, ok := keyVaultsCache[keyVaultUrl]; ok {
		return existing, nil
	}

	list, err := client.VaultsClient.ListComplete(ctx, utils.Int32(1000))
	if err != nil {
		return nil, fmt.Errorf("failed to list Key Vaults %v", err)
	}

	for list.NotDone() {
		v := list.Value()

		if v.ID == nil || *v.ID == "" {
			return nil, fmt.Errorf("v.ID was nil")
		}

		id, err := parse.KeyVaultID(*v.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ID for Key Vault URI %q: %s", *v.ID, err)
		}

		// resp does not appear to contain the vault properties, so lets fetch them
		vault, err := client.VaultsClient.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(vault.Response) {
				if e := list.NextWithContext(ctx); e != nil {
					return nil, fmt.Errorf("failed to get next vault on KeyVault url %q : %+v", keyVaultUrl, err)
				}
				continue
			}
			return nil, fmt.Errorf("failed to make Read request on KeyVault %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		vaultDetails, err := populateKeyVaultDetails(vault)
		if err != nil {
			return nil, err
		}

		keyVaultsCache[*vault.Properties.VaultURI] = vaultDetails

		if e := list.NextWithContext(ctx); e != nil {
			return nil, fmt.Errorf("failed to get next vault on KeyVault url %q : %+v", keyVaultUrl, err)
		}
	}

	if existing, ok := keyVaultsCache[keyVaultUrl]; ok {
		return existing, nil
	}

	return nil, nil
}

func populateKeyVaultDetails(props keyvault.Vault) (*keyVaultDetails, error) {
	if props.ID == nil || *props.ID == "" {
		return nil, fmt.Errorf("`id` was nil or empty for Account %q", props.Name)
	}

	id, err := parse.KeyVaultID(*props.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as key vault ID: %+v", *props.ID, err)
	}

	if props.Properties == nil || props.Properties.VaultURI == nil {
		return nil, fmt.Errorf("KeyVault %q (Resource Group %q) has nil ID, properties or vault URI", id.Name, id.ResourceGroup)
	}

	return &keyVaultDetails{
		ID:            *props.ID,
		Name:          id.Name,
		ResourceGroup: id.ResourceGroup,
		Properties:    props.Properties,
		Tags:          props.Tags,
	}, nil
}
