package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type ClusterId struct {
	ResourceGroup      string
	ManagedClusterName string
}

func ClusterID(input string) (*ClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cluster := ClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if cluster.ManagedClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cluster, nil
}
