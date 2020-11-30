package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ResourceGroupId struct {
	SubscriptionId string
	ResourceGroup  string
}

func NewResourceGroupID(subscriptionId, resourceGroup string) ResourceGroupId {
	return ResourceGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
	}
}

func (id ResourceGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup)
}

// ResourceGroupID parses a ResourceGroup ID into an ResourceGroupId struct
func ResourceGroupID(input string) (*ResourceGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ResourceGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
