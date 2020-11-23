package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitalTwinsId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewDigitalTwinsID(subscriptionId, resourceGroup string, name string) DigitalTwinsId {
	return DigitalTwinsId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id DigitalTwinsId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func DigitalTwinsID(input string) (*DigitalTwinsId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DigitalTwins ID %q: %+v", input, err)
	}

	digitalTwins := DigitalTwinsId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if digitalTwins.Name, err = id.PopSegment("digitalTwinsInstances"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &digitalTwins, nil
}
