package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MaintenanceConfigurationId struct {
	ResourceGroup string
	Name          string
}

func MaintenanceConfigurationID(input string) (*MaintenanceConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse MaintenanceConfiguration ID %q: %+v", input, err)
	}

	maintenanceConfiguration := MaintenanceConfigurationId{
		ResourceGroup: id.ResourceGroup,
	}

	if maintenanceConfiguration.Name, err = id.PopSegment("maintenanceconfigurations"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &maintenanceConfiguration, nil
}
