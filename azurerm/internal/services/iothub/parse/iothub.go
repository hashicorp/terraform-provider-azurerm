package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IotHubId struct {
	Name          string
	ResourceGroup string
}

func IotHubID(input string) (*IotHubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	iothub := IotHubId{
		ResourceGroup: id.ResourceGroup,
	}

	if iothub.Name, err = id.PopSegment("IotHubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &iothub, nil
}
