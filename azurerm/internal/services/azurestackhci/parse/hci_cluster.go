package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HCIClusterId struct {
	ResourceGroup string
	Name          string
}

func HCIClusterID(input string) (*HCIClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing hciCluster ID %q: %+v", input, err)
	}

	hciCluster := HCIClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if hciCluster.Name, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &hciCluster, nil
}
