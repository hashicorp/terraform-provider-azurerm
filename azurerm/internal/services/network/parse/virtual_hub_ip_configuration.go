package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubIpConfigurationId struct {
	SubscriptionId      string // placeholder for the generator
	ResourceGroup       string
	VirtualHubName      string
	IpConfigurationName string
}

func VirtualHubIpConfigurationID(input string) (*VirtualHubIpConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing virtualHubIP ID %q: %+v", input, err)
	}

	virtualHubIP := VirtualHubIpConfigurationId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHubIP.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if virtualHubIP.IpConfigurationName, err = id.PopSegment("ipConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHubIP, nil
}
