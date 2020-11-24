package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitalTwinsInstanceId struct {
	ResourceGroup string
	Name          string
}

func NewDigitalTwinsInstanceID(resourcegroup string, name string) DigitalTwinsInstanceId {
	return DigitalTwinsInstanceId{
		ResourceGroup: resourcegroup,
		Name:          name,
	}
}

func (id DigitalTwinsInstanceId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DigitalTwinsInstanceID(input string) (*DigitalTwinsInstanceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DigitalTwins ID %q: %+v", input, err)
	}

	digitalTwins := DigitalTwinsInstanceId{
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
