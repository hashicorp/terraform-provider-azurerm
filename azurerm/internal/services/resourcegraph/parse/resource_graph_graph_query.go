package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGraphGraphQueryId struct {
	ResourceGroup string
	ResourceName  string
}

func ResourceGraphGraphQueryID(input string) (*ResourceGraphGraphQueryId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse ResourceGraphGraphQuery ID %q: %+v", input, err)
	}

	resourceGraphGraphQuery := ResourceGraphGraphQueryId{
		ResourceGroup: id.ResourceGroup,
	}
	if resourceGraphGraphQuery.ResourceName, err = id.PopSegment("queries"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceGraphGraphQuery, nil
}
