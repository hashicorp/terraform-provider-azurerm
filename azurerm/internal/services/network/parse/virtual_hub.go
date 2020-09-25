package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubId struct {
	ResourceGroup string
	Name          string
}

func NewVirtualHubID(resourceGroup, name string) VirtualHubId {
	return VirtualHubId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}
func (id VirtualHubId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func VirtualHubID(input string) (*VirtualHubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Hub ID %q: %+v", input, err)
	}

	virtualHub := VirtualHubId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHub.Name, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHub, nil
}
