package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type ClusterId struct {
	ResourceGroup string
	Name          string
}

func ClusterID(input string) (*ClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cluster := ClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if cluster.Name, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cluster, nil
}
