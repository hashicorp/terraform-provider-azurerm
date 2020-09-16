package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseSparkPoolId struct {
	Workspace *SynapseWorkspaceId
	Name      string
}

func SynapseSparkPoolID(input string) (*SynapseSparkPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapse Spark Pool ID %q: %+v", input, err)
	}

	synapseSparkPool := SynapseSparkPoolId{
		Workspace: &SynapseWorkspaceId{
			SubscriptionID: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		},
	}
	if synapseSparkPool.Workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if synapseSparkPool.Name, err = id.PopSegment("bigDataPools"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &synapseSparkPool, nil
}
