package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CdnEndpointId struct {
	ResourceGroup string
	ProfileName   string
	Name          string
}

func CdnEndpointID(input string) (*CdnEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Endpoint ID %q: %+v", input, err)
	}

	endpoint := CdnEndpointId{
		ResourceGroup: id.ResourceGroup,
	}

	if endpoint.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}

	if endpoint.Name, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &endpoint, nil
}
