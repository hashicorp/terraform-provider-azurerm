package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NatGatewayId struct {
	Name          string
	ResourceGroup string
}

func NatGatewayID(input string) (*NatGatewayId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	natGateway := NatGatewayId{
		ResourceGroup: id.ResourceGroup,
	}

	if natGateway.Name, err = id.PopSegment("natGateways"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &natGateway, nil
}
