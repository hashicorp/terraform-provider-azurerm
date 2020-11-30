package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateDnsZoneId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewPrivateDnsZoneID(subscriptionId, resourceGroup, name string) PrivateDnsZoneId {
	return PrivateDnsZoneId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id PrivateDnsZoneId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// PrivateDnsZoneID parses a PrivateDnsZone ID into an PrivateDnsZoneId struct
func PrivateDnsZoneID(input string) (*PrivateDnsZoneId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateDnsZoneId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("privateDnsZones"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
