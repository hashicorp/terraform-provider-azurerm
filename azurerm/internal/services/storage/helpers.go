package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func (client Client) FindResourceGroup(ctx context.Context, accountName string) (*string, error) {
	accounts, err := client.AccountsClient.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error listing Storage Accounts (to find Resource Group for %q): %s", accountName, err)
	}

	if accounts.Value == nil {
		return nil, nil
	}

	var resourceGroup *string
	for _, account := range *accounts.Value {
		if account.Name == nil || account.ID == nil {
			continue
		}

		if strings.EqualFold(accountName, *account.Name) {
			id, err := azure.ParseAzureResourceID(*account.ID)
			if err != nil {
				return nil, fmt.Errorf("Error parsing ID for Storage Account %q: %s", accountName, err)
			}

			resourceGroup = &id.ResourceGroup
			break
		}
	}

	return resourceGroup, nil
}

func (client Client) findAccountKey(ctx context.Context, resourceGroup, accountName string) (*string, error) {
	props, err := client.AccountsClient.ListKeys(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error Listing Keys for Storage Account %q (Resource Group %q): %+v", accountName, resourceGroup, err)
	}

	if props.Keys == nil || len(*props.Keys) == 0 {
		return nil, fmt.Errorf("Keys were nil for Storage Account %q (Resource Group %q): %+v", accountName, resourceGroup, err)
	}

	keys := *props.Keys
	firstKey := keys[0].Value
	return firstKey, nil
}
