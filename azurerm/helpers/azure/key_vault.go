package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
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
			return "", fmt.Errorf("failed to find KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
		}
		return "", fmt.Errorf("failed to make Read request on KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if resp.Properties == nil || resp.Properties.VaultURI == nil {
		return "", fmt.Errorf("vault (%s) response properties or VaultURI is nil", keyVaultId)
	}

	return *resp.Properties.VaultURI, nil
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
		return false, fmt.Errorf("failed to make Read request on KeyVault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if resp.Properties == nil || resp.Properties.VaultURI == nil {
		return false, fmt.Errorf("vault (%s) response properties or VaultURI is nil", keyVaultId)
	}

	return true, nil
}
