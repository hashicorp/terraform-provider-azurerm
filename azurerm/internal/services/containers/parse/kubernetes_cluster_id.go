package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type KubernetesClusterId struct {
	Name          string
	ResourceGroup string
}

func KubernetesClusterID(input string) (*KubernetesClusterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	cluster := KubernetesClusterId{
		ResourceGroup: id.ResourceGroup,
	}

	if cluster.Name, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cluster, nil
}
