package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KubernetesNodePoolId struct {
	Name          string
	ClusterName   string
	ResourceGroup string
}

func KubernetesNodePoolID(input string) (*KubernetesNodePoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	pool := KubernetesNodePoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if pool.ClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}

	if pool.Name, err = id.PopSegment("agentPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &pool, nil
}
