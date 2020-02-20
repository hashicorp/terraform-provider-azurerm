package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabaseMigrationServiceId struct {
	ResourceGroup string
	Name          string
}

func DatabaseMigrationServiceID(input string) (*DatabaseMigrationServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Database Migration Service ID %q: %+v", input, err)
	}

	server := DatabaseMigrationServiceId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("services"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
