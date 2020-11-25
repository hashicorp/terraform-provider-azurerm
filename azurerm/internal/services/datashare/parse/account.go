package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AccountId struct {
	ResourceGroup string
	Name          string
}

func (id AccountId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func NewAccountId(resourceGroup, name string) AccountId {
	return AccountId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func AccountID(input string) (*AccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DataShareAccount ID %q: %+v", input, err)
	}

	dataShareAccount := AccountId{
		ResourceGroup: id.ResourceGroup,
	}
	if dataShareAccount.Name, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &dataShareAccount, nil
}
