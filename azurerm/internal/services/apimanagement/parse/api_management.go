package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiManagementId struct {
	ResourceGroup string
	ServiceName   string
}

func ApiManagementID(input string) (*ApiManagementId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	apiManagement := ApiManagementId{
		ResourceGroup: id.ResourceGroup,
	}

	if apiManagement.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &apiManagement, nil
}
