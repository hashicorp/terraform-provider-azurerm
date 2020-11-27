package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateEndpointId struct {
	ResourceGroup string
	Name          string
}

func PrivateEndpointID(input string) (*PrivateEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Private Endpoint ID %q: %+v", input, err)
	}

	privateEndpoint := PrivateEndpointId{
		ResourceGroup: id.ResourceGroup,
	}

	if privateEndpoint.Name, err = id.PopSegment("privateEndpoints"); err != nil {
		return nil, err
	}

	return &privateEndpoint, nil
}
