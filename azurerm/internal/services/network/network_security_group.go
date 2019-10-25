package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkSecurityGroupResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseNetworkSecurityGroupResourceID(input string) (*NetworkSecurityGroupResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Network Security Group ID %q: %+v", input, err)
	}

	networkSecurityGroup := NetworkSecurityGroupResourceID{
		Base: *id,
		Name: id.Path["networkSecurityGroups"],
	}

	if networkSecurityGroup.Name == "" {
		return nil, fmt.Errorf("ID was missing the `networkSecurityGroups` element")
	}

	return &networkSecurityGroup, nil
}
