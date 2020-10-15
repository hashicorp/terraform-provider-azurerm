package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubIPConfigurationId struct {
	ResourceGroup  string
	VirtualHubName string
	Name           string
}

func VirtualHubIPConfigurationID(input string) (*VirtualHubIPConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing virtualHubIPConfiguration ID %q: %+v", input, err)
	}

	virtualHubIpConfiguration := VirtualHubIPConfigurationId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHubIpConfiguration.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if virtualHubIpConfiguration.Name, err = id.PopSegment("ipConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHubIpConfiguration, nil
}
