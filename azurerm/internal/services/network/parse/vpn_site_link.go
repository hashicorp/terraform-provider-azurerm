package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VpnSiteLinkId struct {
	ResourceGroup string
	Site          string
	Name          string
}

func (id VpnSiteLinkId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnSites/%s/vpnSiteLinks/%s",
		subscriptionId, id.ResourceGroup, id.Site, id.Name)
}

func NewVpnSiteLinkID(vpnSiteId VpnSiteId, name string) VpnSiteLinkId {
	return VpnSiteLinkId{
		ResourceGroup: vpnSiteId.ResourceGroup,
		Site:          vpnSiteId.Name,
		Name:          name,
	}
}

func VpnSiteLinkID(input string) (*VpnSiteLinkId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Vpn Site Link ID %q: %+v", input, err)
	}

	vpnSiteLinkId := VpnSiteLinkId{
		ResourceGroup: id.ResourceGroup,
	}

	if vpnSiteLinkId.Site, err = id.PopSegment("vpnSites"); err != nil {
		return nil, err
	}

	if vpnSiteLinkId.Name, err = id.PopSegment("vpnSiteLinks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vpnSiteLinkId, nil
}
