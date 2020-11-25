package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ShareId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func (id ShareId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataShare/accounts/%s/shares/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.AccountName, id.Name)
}

func NewShareId(resourceGroup, accountName, name string) ShareId {
	return ShareId{
		ResourceGroup: resourceGroup,
		AccountName:   accountName,
		Name:          name,
	}
}

func ShareID(input string) (*ShareId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DataShare ID %q: %+v", input, err)
	}

	DataShare := ShareId{
		ResourceGroup: id.ResourceGroup,
	}
	if DataShare.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if DataShare.Name, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &DataShare, nil
}
