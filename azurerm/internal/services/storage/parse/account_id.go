package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountID struct {
	Name           string
	ResourceGroup  string
	SubscriptionId string
}

func ParseAccountID(input string) (*AccountID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := AccountID{
		ResourceGroup:  id.ResourceGroup,
		SubscriptionId: id.SubscriptionID,
	}

	if account.Name, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}
