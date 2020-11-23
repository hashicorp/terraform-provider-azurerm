package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitalTwinsId struct {
	ResourceGroup string
	Name          string
}

func NewDigitalTwinsID(resourcegroup string, name string) DigitalTwinsId {
	return DigitalTwinsId{
		ResourceGroup: resourcegroup,
		Name:          name,
	}
}

func (id DigitalTwinsId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DigitalTwinsID(input string) (*DigitalTwinsId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DigitalTwins ID %q: %+v", input, err)
	}

	digitalTwins := DigitalTwinsId{
		ResourceGroup: id.ResourceGroup,
	}
	if digitalTwins.Name, err = id.PopSegment("digitalTwinsInstances"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &digitalTwins, nil
}
