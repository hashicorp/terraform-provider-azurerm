package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitalTwinsInstanceId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewDigitalTwinsInstanceID(subscriptionId, resourceGroup, name string) DigitalTwinsInstanceId {
	return DigitalTwinsInstanceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id DigitalTwinsInstanceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func DigitalTwinsInstanceID(input string) (*DigitalTwinsInstanceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing DigitalTwins ID %q: %+v", input, err)
	}

	digitalTwinsInstance := DigitalTwinsInstanceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if digitalTwinsInstance.Name, err = id.PopSegment("digitalTwinsInstances"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &digitalTwinsInstance, nil
}
