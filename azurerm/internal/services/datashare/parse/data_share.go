package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataShareAccountId struct {
	ResourceGroup string
	Name          string
}

type DataShareId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func DataShareAccountID(input string) (*DataShareAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DataShareAccount ID %q: %+v", input, err)
	}

	dataShareAccount := DataShareAccountId{
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

func DataShareID(input string) (*DataShareId, error) {
	var id, err = azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DataShare ID %q: %+v", input, err)
	}

	DataShare := DataShareId{
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
