package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type DevTestVirtualNetworkId struct {
	SubscriptionId     string
	ResourceGroup      string
	LabName            string
	VirtualNetworkName string
}

func NewDevTestVirtualNetworkID(subscriptionId, resourceGroup, labName, virtualNetworkName string) DevTestVirtualNetworkId {
	return DevTestVirtualNetworkId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		LabName:            labName,
		VirtualNetworkName: virtualNetworkName,
	}
}

func (id DevTestVirtualNetworkId) String() string {
	segments := []string{
		fmt.Sprintf("Virtual Network Name %q", id.VirtualNetworkName),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dev Test Virtual Network", segmentsStr)
}

func (id DevTestVirtualNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/virtualNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.VirtualNetworkName)
}

// DevTestVirtualNetworkID parses a DevTestVirtualNetwork ID into an DevTestVirtualNetworkId struct
func DevTestVirtualNetworkID(input string) (*DevTestVirtualNetworkId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DevTestVirtualNetworkId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LabName, err = id.PopSegment("labs"); err != nil {
		return nil, err
	}
	if resourceId.VirtualNetworkName, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
