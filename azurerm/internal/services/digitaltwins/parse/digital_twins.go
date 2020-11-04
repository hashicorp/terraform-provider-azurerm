package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitaltwinsDigitalTwinId struct {
	ResourceGroup string
	Name          string
}

func NewDigitaltwinsDigitalTwinID(resourcegroup string, name string) DigitaltwinsDigitalTwinId {
	return DigitaltwinsDigitalTwinId{
		ResourceGroup: resourcegroup,
		Name:          name,
	}
}

func (id DigitaltwinsDigitalTwinId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DigitaltwinsDigitalTwinID(input string) (*DigitaltwinsDigitalTwinId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing digitaltwinsDigitalTwin ID %q: %+v", input, err)
	}

	digitaltwinsDigitalTwin := DigitaltwinsDigitalTwinId{
		ResourceGroup: id.ResourceGroup,
	}
	if digitaltwinsDigitalTwin.Name, err = id.PopSegment("digitalTwinsInstances"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &digitaltwinsDigitalTwin, nil
}
