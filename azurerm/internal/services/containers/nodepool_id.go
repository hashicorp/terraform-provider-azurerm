package containers

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KubernetesNodePoolID struct {
	Name          string
	ClusterName   string
	ResourceGroup string
}

func ParseKubernetesNodePoolID(input string) (*KubernetesNodePoolID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	pool := KubernetesNodePoolID{
		ResourceGroup: id.ResourceGroup,
	}

	pool.ClusterName, err = id.PopSegment("managedClusters")
	if err != nil {
		return nil, err
	}

	pool.Name, err = id.PopSegment("agentPools")
	if err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &pool, nil
}
