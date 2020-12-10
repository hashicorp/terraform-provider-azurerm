package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkSwiftConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
	ConfigName     string
}

func NewVirtualNetworkSwiftConnectionID(subscriptionId, resourceGroup, siteName, configName string) VirtualNetworkSwiftConnectionId {
	return VirtualNetworkSwiftConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
		ConfigName:     configName,
	}
}

func (id VirtualNetworkSwiftConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Config Name %q", id.ConfigName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Swift Connection", segmentsStr)
}

func (id VirtualNetworkSwiftConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/config/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.ConfigName)
}

// VirtualNetworkSwiftConnectionID parses a VirtualNetworkSwiftConnection ID into an VirtualNetworkSwiftConnectionId struct
func VirtualNetworkSwiftConnectionID(input string) (*VirtualNetworkSwiftConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkSwiftConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.ConfigName, err = id.PopSegment("config"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
