package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitalTwinId struct {
	ResourceGroup string
	Name          string
}

func NewDigitalTwinID(resourcegroup string, name string) DigitalTwinId {
	return DigitalTwinId{
		ResourceGroup: resourcegroup,
		Name:          name,
	}
}

func (id DigitalTwinId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DigitalTwinID(input string) (*DigitalTwinId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DigitalTwin ID %q: %+v", input, err)
	}

	digitalTwins := DigitalTwinId{
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
