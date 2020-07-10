package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiManagementLoggerId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func ApiManagementLoggerID(input string) (*ApiManagementLoggerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	logger := ApiManagementLoggerId{
		ResourceGroup: id.ResourceGroup,
	}

	if logger.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if logger.Name, err = id.PopSegment("loggers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &logger, nil
}
