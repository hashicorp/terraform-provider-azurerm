package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseBigDataPoolId struct {
	Workspace *SynapseWorkspaceId
	Name      string
}

func SynapseBigDataPoolID(input string) (*SynapseBigDataPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapseBigDataPool ID %q: %+v", input, err)
	}

	synapseBigDataPool := SynapseBigDataPoolId{
		Workspace: &SynapseWorkspaceId{
			SubscriptionID: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		},
	}
	if synapseBigDataPool.Workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if synapseBigDataPool.Name, err = id.PopSegment("bigDataPools"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &synapseBigDataPool, nil
}
