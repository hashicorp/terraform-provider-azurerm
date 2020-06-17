package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubConnectionResourceID struct {
	ResourceGroup  string
	VirtualHubId   string
	VirtualHubName string
	Name           string
}

func ParseVirtualHubConnectionID(input string) (*VirtualHubConnectionResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Hub Connection ID %q: %+v", input, err)
	}

	// /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tom-dev99/providers/Microsoft.Network/virtualHubs/tom-devvh/hubVirtualNetworkConnections/first
	connection := VirtualHubConnectionResourceID{
		ResourceGroup:  id.ResourceGroup,
		VirtualHubName: id.Path["virtualHubs"],
		Name:           id.Path["hubVirtualNetworkConnections"],
	}

	connectionIndex := strings.Index(input, "/hubVirtualNetworkConnections")
	if connectionIndex == -1 {
		return nil, fmt.Errorf("parsing virtual hub resource id from hub virtual network connection resource id")
	} else {
		connection.VirtualHubId = input[:connectionIndex]
	}

	if connection.VirtualHubName == "" {
		return nil, fmt.Errorf("ID was missing the `virtualHubs` element")
	}

	if connection.Name == "" {
		return nil, fmt.Errorf("ID was missing the `hubVirtualNetworkConnections` element")
	}

	return &connection, nil
}
