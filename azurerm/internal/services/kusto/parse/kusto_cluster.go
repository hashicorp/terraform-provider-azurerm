package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KustoClusterId struct {
	ResourceGroup string
	Name          string
}

func KustoClusterID(input string) (*KustoClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Cluster ID %q: %+v", input, err)
	}

	cluster := KustoClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if cluster.Name, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cluster, nil
}
