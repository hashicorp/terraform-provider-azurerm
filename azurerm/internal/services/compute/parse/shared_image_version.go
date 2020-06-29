package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SharedImageVersionId struct {
	ResourceGroup string
	Version       string
	Gallery       string
	Name          string
}

func SharedImageVersionID(input string) (*SharedImageVersionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Image ID %q: %+v", input, err)
	}

	set := SharedImageVersionId{
		ResourceGroup: id.ResourceGroup,
	}

	if set.Gallery, err = id.PopSegment("galleries"); err != nil {
		return nil, err
	}

	if set.Name, err = id.PopSegment("images"); err != nil {
		return nil, err
	}

	if set.Version, err = id.PopSegment("versions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &set, nil
}
