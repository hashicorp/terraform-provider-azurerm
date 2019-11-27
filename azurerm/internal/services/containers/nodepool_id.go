package containers

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KubernetesNodePoolID struct {
	Name          string
	ClusterName   string
	ResourceGroup string

	ID azure.ResourceID
}

func ParseKubernetesNodePoolID(id string) (*KubernetesNodePoolID, error) {
	clusterId, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}

	resourceGroup := clusterId.ResourceGroup
	if resourceGroup == "" {
		return nil, fmt.Errorf("%q is missing a Resource Group", id)
	}

	clusterName := clusterId.Path["managedClusters"]
	if clusterName == "" {
		return nil, fmt.Errorf("%q is missing the `managedClusters` segment", id)
	}

	nodePoolName := clusterId.Path["agentPools"]
	if nodePoolName == "" {
		return nil, fmt.Errorf("%q is missing the `agentPools` segment", id)
	}

	output := KubernetesNodePoolID{
		Name:          nodePoolName,
		ClusterName:   clusterName,
		ResourceGroup: resourceGroup,
		ID:            *clusterId,
	}
	return &output, nil
}
