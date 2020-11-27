package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MediaStreamingEndpointId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func MediaStreamingEndpointID(input string) (*MediaStreamingEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Streaming Endpoint ID %q: %+v", input, err)
	}

	endpoint := MediaStreamingEndpointId{
		ResourceGroup: id.ResourceGroup,
	}

	if endpoint.AccountName, err = id.PopSegment("mediaservices"); err != nil {
		return nil, err
	}

	if endpoint.Name, err = id.PopSegment("streamingendpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &endpoint, nil
}
