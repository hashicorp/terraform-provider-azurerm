package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KustoTableId struct {
	ResourceGroup string
	Cluster       string
	Database      string
	Name          string
}

func KustoTableID(input string) (*KustoTableId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Table ID %q: %+v", input, err)
	}

	table := KustoTableId{
		ResourceGroup: id.ResourceGroup,
	}

	if table.Cluster, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if table.Database, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}

	if table.Name, err = id.PopSegment("Tables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &table, nil
}
