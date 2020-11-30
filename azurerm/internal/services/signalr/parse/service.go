package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceId struct {
	ResourceGroup string
	SignalRName   string
}

func ServiceID(input string) (*ServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse SignalR Service ID %q: %+v", input, err)
	}

	service := ServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.SignalRName, err = id.PopSegment("SignalR"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
