package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkSwiftConnectionId struct {
	ResourceGroup string
	SiteName      string
}

func VirtualNetworkSwiftConnectionID(resourceId string) (*VirtualNetworkSwiftConnectionId, error) {
	id, err := azure.ParseAzureResourceID(resourceId)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}

	virtualNetworkId := &VirtualNetworkSwiftConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualNetworkId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	return virtualNetworkId, nil
}
