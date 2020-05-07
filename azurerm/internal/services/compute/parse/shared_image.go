package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageId struct {
	ResourceGroup string
	Gallery       string
	Name          string
}

func SharedImageID(input string) (*SharedImageId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Image ID %q: %+v", input, err)
	}

	set := SharedImageId{
		ResourceGroup: id.ResourceGroup,
	}

	if set.Gallery, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}

	if set.Name, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
