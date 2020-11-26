package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProximityPlacementGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewProximityPlacementGroupID(subscriptionId, resourceGroup, name string) ProximityPlacementGroupId {
	return ProximityPlacementGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ProximityPlacementGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/proximityPlacementGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ProximityPlacementGroupID parses a ProximityPlacementGroup ID into an ProximityPlacementGroupId struct
func ProximityPlacementGroupID(input string) (*ProximityPlacementGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ProximityPlacementGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("proximityPlacementGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
