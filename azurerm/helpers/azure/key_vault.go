package azure

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2018-02-14/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func GetKeyVaultNameFromBaseUrl(keyVaultUrl string) (string, error) {
	groups := regexp.MustCompile(`^https://(.+)\.vault\.azure\.net/?$`).FindStringSubmatch(keyVaultUrl)
	if len(groups) != 2 {
		return "", fmt.Errorf("parsing keyVaultUrl: %q, expected group: 2", keyVaultUrl)
	}

	return groups[1], nil
}

func GetKeyVaultIDFromBaseUrl(ctx context.Context, client *resources.Client, keyVaultUrl string) (*string, error) {
	name, err := GetKeyVaultNameFromBaseUrl(keyVaultUrl)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("name eq '%s'", name)
	keyVaults, err := client.List(ctx, filter, "", nil)
	if err != nil {
		return nil, fmt.Errorf("listing key vault with name %q: %+v", name, err)
	}
	values := keyVaults.Values()
	if len(values) == 0 {
		return nil, nil
	} else if len(values) > 1 {
		return nil, fmt.Errorf("more than one key Vault found with Url: %q", keyVaultUrl)
	}

	return values[0].ID, nil
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
