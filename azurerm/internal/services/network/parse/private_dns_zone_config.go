package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateDnsZoneConfigId struct {
	SubscriptionId          string
	ResourceGroup           string
	PrivateEndpointName     string
	PrivateDnsZoneGroupName string
	Name                    string
}

func NewPrivateDnsZoneConfigID(subscriptionId, resourceGroup, privateEndpointName, privateDnsZoneGroupName, name string) PrivateDnsZoneConfigId {
	return PrivateDnsZoneConfigId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		PrivateEndpointName:     privateEndpointName,
		PrivateDnsZoneGroupName: privateDnsZoneGroupName,
		Name:                    name,
	}
}

func (id PrivateDnsZoneConfigId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Private Dns Zone Group Name %q", id.PrivateDnsZoneGroupName),
		fmt.Sprintf("Private Endpoint Name %q", id.PrivateEndpointName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Dns Zone Config", segmentsStr)
}

func (id PrivateDnsZoneConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s/privateDnsZoneGroups/%s/privateDnsZoneConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PrivateEndpointName, id.PrivateDnsZoneGroupName, id.Name)
}

// PrivateDnsZoneConfigID parses a PrivateDnsZoneConfig ID into an PrivateDnsZoneConfigId struct
func PrivateDnsZoneConfigID(input string) (*PrivateDnsZoneConfigId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateDnsZoneConfigId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PrivateEndpointName, err = id.PopSegment("privateEndpoints"); err != nil {
		return nil, err
	}
	if resourceId.PrivateDnsZoneGroupName, err = id.PopSegment("privateDnsZoneGroups"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("privateDnsZoneConfigs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
