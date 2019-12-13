package parsers

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountID struct {
	Name          string
	ResourceGroup string

	ID azure.ResourceID
}

func ParseAccountID(id string) (*AccountID, error) {
	storageID, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}

	resourceGroup := storageID.ResourceGroup
	if resourceGroup == "" {
		return nil, fmt.Errorf("%q is missing a Resource Group", id)
	}

	storageAccountName := storageID.Path["storageAccounts"]
	if storageAccountName == "" {
		return nil, fmt.Errorf("%q is missing the `storageAccounts` segment", id)
	}

	accountId := AccountID{
		Name:          storageAccountName,
		ResourceGroup: resourceGroup,
		ID:            *storageID,
	}
	return &accountId, nil
}
