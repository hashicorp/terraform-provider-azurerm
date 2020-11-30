package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateDnsZoneGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewPrivateDnsZoneGroupID(subscriptionId, resourceGroup, name string) PrivateDnsZoneGroupId {
	return PrivateDnsZoneGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id PrivateDnsZoneGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZoneGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// PrivateDnsZoneGroupID parses a PrivateDnsZoneGroup ID into an PrivateDnsZoneGroupId struct
func PrivateDnsZoneGroupID(input string) (*PrivateDnsZoneGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateDnsZoneGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("privateDnsZoneGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
