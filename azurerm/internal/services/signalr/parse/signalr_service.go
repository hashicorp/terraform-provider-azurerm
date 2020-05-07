package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SignalRServiceId struct {
	ResourceGroup string
	Name          string
}

func SignalRServiceID(input string) (*SignalRServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse SignalR Service ID %q: %+v", input, err)
	}

	service := SignalRServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.Name, err = id.PopSegment("SignalR"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
