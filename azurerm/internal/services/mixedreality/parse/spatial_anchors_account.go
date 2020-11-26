package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpatialAnchorsAccountId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewSpatialAnchorsAccountID(subscriptionId, resourceGroup, name string) SpatialAnchorsAccountId {
	return SpatialAnchorsAccountId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id SpatialAnchorsAccountId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// SpatialAnchorsAccountID parses a SpatialAnchorsAccount ID into an SpatialAnchorsAccountId struct
func SpatialAnchorsAccountID(input string) (*SpatialAnchorsAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpatialAnchorsAccountId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("spatialAnchorsAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
