package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ImageId struct {
	ResourceGroup string
	Name          string
}

func NewImageId(resourceGroup, name string) ImageId {
	return ImageId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id ImageId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/images/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func ImageID(input string) (*ImageId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Image ID %q: %+v", input, err)
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
