package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HyperConvergedClusterId struct {
	ResourceGroup string
	Name          string
}

func HyperConvergedClusterID(input string) (*HyperConvergedClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing hyperConvergedCluster ID %q: %+v", input, err)
	}

	hyperConvergedCluster := HyperConvergedClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if hyperConvergedCluster.Name, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &hyperConvergedCluster, nil
}
