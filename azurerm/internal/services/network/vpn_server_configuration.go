package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VpnServerConfigurationResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseVpnServerConfigurationID(input string) (*VpnServerConfigurationResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse VPN Server Configuration ID %q: %+v", input, err)
	}

	vpnServerConfiguration := VpnServerConfigurationResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if vpnServerConfiguration.Name, err = id.PopSegment("vpnServerConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vpnServerConfiguration, nil
}

func ValidateVpnServerConfigurationID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if _, err := ParseVpnServerConfigurationID(v); err != nil {
		return nil, []error{err}
	}

	return nil, nil
}
