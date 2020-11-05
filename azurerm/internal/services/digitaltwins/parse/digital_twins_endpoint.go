package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitaltwinsEndpointId struct {
	ResourceGroup string
	ResourceName  string
	Name          string
}

func NewDigitaltwinsEndpointID(resourcegroup string, resourcename string, name string) DigitaltwinsEndpointId {
	return DigitaltwinsEndpointId{
		ResourceGroup: resourcegroup,
		ResourceName:  resourcename,
		Name:          name,
	}
}

func (id DigitaltwinsEndpointId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s/endpoints/%s", subscriptionId, id.ResourceGroup, id.ResourceName, id.Name)
}

func DigitaltwinsEndpointID(input string) (*DigitaltwinsEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DigitalTwinsEndpoint ID %q: %+v", input, err)
	}

	DigitalTwinsEndpoint := DigitaltwinsEndpointId{
		ResourceGroup: id.ResourceGroup,
	}
	if DigitalTwinsEndpoint.ResourceName, err = id.PopSegment("digitalTwinsInstances"); err != nil {
		return nil, err
	}
	if DigitalTwinsEndpoint.Name, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &DigitalTwinsEndpoint, nil
}
