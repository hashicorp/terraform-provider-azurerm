package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkSecurityGroupResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseNetworkSecurityGroupID(input string) (*NetworkSecurityGroupResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Network Security Group ID %q: %+v", input, err)
	}

	networkSecurityGroup := NetworkSecurityGroupResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if networkSecurityGroup.Name, err = id.PopSegment("networkSecurityGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &networkSecurityGroup, nil
}
