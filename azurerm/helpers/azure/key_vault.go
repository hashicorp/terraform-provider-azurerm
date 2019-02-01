package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func GetKeyVaultBaseUrlFromID(ctx context.Context, client keyvault.VaultsClient, keyVaultId string) (string, error) {

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
		return "", fmt.Errorf("resource id does not contain `valuts`: %q", keyVaultId)
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

func GetKeyVaultIDFromBaseUrl(ctx context.Context, client keyvault.VaultsClient, keyVaultUrl string) (string, error) {

	list, err := client.ListComplete(ctx, utils.Int32(1))
	if err != nil {
		return "", fmt.Errorf("Error GetKeyVaultId unable to list Key Vaults %v", err)
	}

	for list.NotDone() {
		v := list.Value()
		if v.ID == nil {
			log.Printf("[DEBUG]GetKeyVaultId: v.ID was nil, continuing")
			continue
		}

		vid, err := ParseAzureResourceID(*v.ID)
		if err != nil {
			log.Printf("[DEBUG] GetKeyVaultId: unable to parse v.ID (%s): %v", *v.ID, err)
			continue
		}
		resourceGroup := vid.ResourceGroup
		name := vid.Path["vaults"]

		//resp does not appear to contain the vault properties, so lets fech them
		get, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			log.Printf("[DEBUG] GetKeyVaultId: Error making Read request on KeyVault %q (Resource Group %q): %+v", name, resourceGroup, err)
			continue
		}

		if get.ID == nil || get.Properties == nil || get.Properties.VaultURI == nil {
			log.Printf("[DEBUG] GetKeyVaultId: KeyVault %q (Resource Group %q) has nil ID, properties or vault URI", name, resourceGroup)
			continue
		}

		if keyVaultUrl == *get.Properties.VaultURI {
			return *get.ID, nil
		}

		e := list.NextWithContext(ctx)
		if e != nil {
			return "", fmt.Errorf("Error GetKeyVaultId: Error getting next value on KeyVault %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return "", fmt.Errorf("Error GetKeyVaultId unable to find Key Vault with url %q", keyVaultUrl)
}

func KeyVaultExists(ctx context.Context, client keyvault.VaultsClient, keyVaultId string) (bool, error) {

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
		return false, fmt.Errorf("resource id does not contain `valuts`: %q", keyVaultId)
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
