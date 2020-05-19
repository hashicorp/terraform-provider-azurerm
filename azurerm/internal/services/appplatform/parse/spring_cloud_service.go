package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudServiceId struct {
	ResourceGroup string
	Name          string
}

func SpringCloudServiceID(input string) (*SpringCloudServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Spring Cloud Service ID %q: %+v", input, err)
	}

	server := SpringCloudServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
