package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountId struct {
	Name           string
	ResourceGroup  string
	SubscriptionId string
}

func NewAccountId(subscriptionId, resourceGroup, name string) AccountId {
	return AccountId{
		Name:           name,
		ResourceGroup:  resourceGroup,
		SubscriptionId: subscriptionId,
	}
}

// the subscriptionId isn't used here, this is just to comply with the interface for now..
func (id AccountId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func AccountID(input string) (*AccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := AccountId{
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
