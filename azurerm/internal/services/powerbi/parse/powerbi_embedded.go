package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PowerBIEmbeddedId struct {
	ResourceGroup string
	Name          string
}

func PowerBIEmbeddedID(input string) (*PowerBIEmbeddedId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse PowerBI Embedded ID %q: %+v", input, err)
	}

	powerbiEmbedded := PowerBIEmbeddedId{
		ResourceGroup: id.ResourceGroup,
	}

	if powerbiEmbedded.Name, err = id.PopSegment("capacities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &powerbiEmbedded, nil
}
