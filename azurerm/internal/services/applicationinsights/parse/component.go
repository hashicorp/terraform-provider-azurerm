package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ComponentId struct {
	SubscriptionId string
	ResourceGroup  string
	ComponentName  string
}

func NewComponentID(subscriptionId, resourceGroup, componentName string) ComponentId {
	return ComponentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ComponentName:  componentName,
	}
}

func (id ComponentId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ComponentName)
}

func ComponentID(input string) (*ComponentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ComponentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.ComponentName, err = id.PopSegment("components"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
