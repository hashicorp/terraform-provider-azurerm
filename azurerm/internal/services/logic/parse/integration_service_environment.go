package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type IntegrationServiceEnvironmentId struct {
	ResourceGroup string
	Name          string
}

func IntegrationServiceEnvironmentID(input string) (*IntegrationServiceEnvironmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Integration Service Environment ID %q: %+v", input, err)
	}

	IntegrationServiceEnvironment := IntegrationServiceEnvironmentId{
		ResourceGroup: id.ResourceGroup,
	}
	if IntegrationServiceEnvironment.Name, err = id.PopSegment("integrationServiceEnvironments"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &IntegrationServiceEnvironment, nil
}
