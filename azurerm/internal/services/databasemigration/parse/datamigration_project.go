package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabaseMigrationProjectId struct {
	ResourceGroup string
	Service       string
	Name          string
}

func DatabaseMigrationProjectID(input string) (*DatabaseMigrationProjectId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Database Migration Project ID %q: %+v", input, err)
	}

	project := DatabaseMigrationProjectId{
		ResourceGroup: id.ResourceGroup,
	}

	if project.Service, err = id.PopSegment("services"); err != nil {
		return nil, err
	}

	if project.Name, err = id.PopSegment("projects"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &project, nil
}
