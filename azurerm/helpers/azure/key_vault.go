package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func GetKeyVaultBaseUrlFromID(ctx context.Context, client *keyvault.VaultsClient, keyVaultId string) (string, error) {
	if keyVaultId == "" {
		return "", fmt.Errorf("keyVaultId is empty")
	}

	id, err := ParseAzureResourceID(keyVaultId)
	if err != nil {
		return "", err
	}
	resourceGroup := id.ResourceGroup

	vaultName, ok := id.Path["vaults"]
	if !ok {
		return "", fmt.Errorf("resource id does not contain `vaults`: %q", keyVaultId)
	}

	resp, err := client.Get(ctx, resourceGroup, vaultName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return "", fmt.Errorf("Error unable to find KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
		}
		return "", fmt.Errorf("Error making Read request on KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if resp.Properties == nil || resp.Properties.VaultURI == nil {
		return "", fmt.Errorf("vault (%s) response properties or VaultURI is nil", keyVaultId)
	}

	return *resp.Properties.VaultURI, nil
}

func GetKeyVaultIDFromBaseUrl(ctx context.Context, client *keyvault.VaultsClient, keyVaultUrl string) (*string, error) {
	list, err := client.ListComplete(ctx, utils.Int32(1000))
	if err != nil {
		return nil, fmt.Errorf("Error GetKeyVaultId unable to list Key Vaults %v", err)
	}

	for list.NotDone() {
		v := list.Value()

		if v.ID == nil {
			return nil, fmt.Errorf("v.ID was nil")
		}

		vid, err := ParseAzureResourceID(*v.ID)
		if err != nil {
			return nil, fmt.Errorf("Error parsing ID for Key Vault URI %q: %s", *v.ID, err)
		}
		resourceGroup := vid.ResourceGroup
		name := vid.Path["vaults"]

		//resp does not appear to contain the vault properties, so lets fetch them
		get, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(get.Response) {
				if e := list.NextWithContext(ctx); e != nil {
					return nil, fmt.Errorf("Error getting next vault on KeyVault url %q : %+v", keyVaultUrl, err)
				}
				continue
			}
			return nil, fmt.Errorf("Error making Read request on KeyVault %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if get.ID == nil || get.Properties == nil || get.Properties.VaultURI == nil {
			return nil, fmt.Errorf("KeyVault %q (Resource Group %q) has nil ID, properties or vault URI", name, resourceGroup)
		}

		if keyVaultUrl == *get.Properties.VaultURI {
			return get.ID, nil
		}

		if e := list.NextWithContext(ctx); e != nil {
			return nil, fmt.Errorf("Error getting next vault on KeyVault url %q : %+v", keyVaultUrl, err)
		}
	}

	// we haven't found it, but Data Sources and Resources need to handle this error separately
	return nil, nil
}

func KeyVaultExists(ctx context.Context, client *keyvault.VaultsClient, keyVaultId string) (bool, error) {
	if keyVaultId == "" {
		return false, fmt.Errorf("keyVaultId is empty")
	}

	id, err := ParseAzureResourceID(keyVaultId)
	if err != nil {
		return false, err
	}
	resourceGroup := id.ResourceGroup

	vaultName, ok := id.Path["vaults"]
	if !ok {
		return false, fmt.Errorf("resource id does not contain `vaults`: %q", keyVaultId)
	}

	resp, err := client.Get(ctx, resourceGroup, vaultName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return false, nil
		}
		return false, fmt.Errorf("Error making Read request on KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if resp.Properties == nil || resp.Properties.VaultURI == nil {
		return false, fmt.Errorf("vault (%s) response properties or VaultURI is nil", keyVaultId)
	}

	return true, nil
}

func KeyVaultCustomizeDiff(d *schema.ResourceDiff, _ interface{}) error {
	if d.HasChange("soft_delete_enabled") {
		if old, new := d.GetChange("soft_delete_enabled"); old.(bool) && !new.(bool) {
			return fmt.Errorf("the property 'soft_delete_enabled' cannot be set to false, enabling the soft delete for a vault is an irreversible action")
		}
	}

	if d.HasChange("purge_protection_enabled") {
		if old, new := d.GetChange("purge_protection_enabled"); old.(bool) && !new.(bool) {
			return fmt.Errorf("the property 'purge_protection_enabled' cannot be set to false, enabling the purge protection for a vault is an irreversible action")
		}
	}

	return nil
}

func KeyVaultGetSoftDeletedState(ctx context.Context, client *keyvault.VaultsClient, name string, location string) (deleteDate interface{}, purgeDate interface{}, err error) {
	softDel, err := client.GetDeleted(ctx, name, location)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get soft delete state information: %+v", err)
	}

	// the logic is this way because the GetDeleted call will return an existing key vault
	// that is not soft deleted, but the Deleted Vault properties will be nil
	if props := softDel.Properties; props != nil {
		var delDate interface{}
		var purgeDate interface{}

		if dd := props.DeletionDate; dd != nil {
			delDate = (*dd).Format(time.RFC3339)
		}
		if pg := props.ScheduledPurgeDate; pg != nil {
			purgeDate = (*pg).Format(time.RFC3339)
		}
		if delDate != nil && purgeDate != nil {
			return delDate, purgeDate, nil
		}
	}

	// this means we found an existing key vault that is not soft deleted
	return nil, nil, nil
}
