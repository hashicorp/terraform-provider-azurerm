package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabaseId struct {
	ResourceGroup string
	Cluster       string
	Name          string
}

func DatabaseID(input string) (*DatabaseId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Database ID %q: %+v", input, err)
	}

	database := DatabaseId{
		ResourceGroup: id.ResourceGroup,
	}

	if database.Cluster, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if database.Name, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &database, nil
}
