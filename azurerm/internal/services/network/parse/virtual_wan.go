package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualWanId struct {
	ResourceGroup string
	Name          string
}

func (id VirtualWanId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualWans/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewVirtualWanID(resourceGroup, name string) VirtualWanId {
	return VirtualWanId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func VirtualWanID(input string) (*VirtualWanId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Virtual Wan ID %q: %+v", input, err)
	}

	vwanId := VirtualWanId{
		ResourceGroup: id.ResourceGroup,
	}

	if vwanId.Name, err = id.PopSegment("virtualWans"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vwanId, nil
}
