package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SparkPoolId struct {
	Workspace       *SynapseWorkspaceId
	BigDataPoolName string
}

func SparkPoolID(input string) (*SparkPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapse Spark Pool ID %q: %+v", input, err)
	}

	synapseSparkPool := SparkPoolId{
		Workspace: &SynapseWorkspaceId{
			SubscriptionID: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		},
	}
	if synapseSparkPool.Workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if synapseSparkPool.BigDataPoolName, err = id.PopSegment("bigDataPools"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &synapseSparkPool, nil
}
