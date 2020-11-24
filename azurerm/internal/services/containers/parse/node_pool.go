package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NodePoolId struct {
	ResourceGroup      string
	ManagedClusterName string
	AgentPoolName      string
}

func NodePoolID(input string) (*NodePoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	pool := NodePoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if pool.ManagedClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}

	if pool.AgentPoolName, err = id.PopSegment("agentPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &pool, nil
}
