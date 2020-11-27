package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EmbeddedId struct {
	ResourceGroup string
	CapacityName  string
}

func EmbeddedID(input string) (*EmbeddedId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse PowerBI Embedded ID %q: %+v", input, err)
	}

	powerbiEmbedded := EmbeddedId{
		ResourceGroup: id.ResourceGroup,
	}

	if powerbiEmbedded.CapacityName, err = id.PopSegment("capacities"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &powerbiEmbedded, nil
}
