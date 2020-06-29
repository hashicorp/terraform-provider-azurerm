package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ImageId struct {
	ResourceGroup string
	Name          string
}

func ImageID(input string) (*ImageId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Image ID %q: %+v", input, err)
	}

	set := ImageId{
		ResourceGroup: id.ResourceGroup,
	}

	if set.Name, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
