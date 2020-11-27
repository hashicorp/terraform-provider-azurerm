package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VpnSiteId struct {
	ResourceGroup string
	Name          string
}

func (id VpnSiteId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnSites/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewVpnSiteID(resourceGroup, name string) VpnSiteId {
	return VpnSiteId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func VpnSiteID(input string) (*VpnSiteId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Vpn Site ID %q: %+v", input, err)
	}

	vpnSiteId := VpnSiteId{
		ResourceGroup: id.ResourceGroup,
	}

	if vpnSiteId.Name, err = id.PopSegment("vpnSites"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vpnSiteId, nil
}
