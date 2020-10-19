package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseSqlPoolId struct {
	Workspace *SynapseWorkspaceId
	Name      string
}

func SynapseSqlPoolID(input string) (*SynapseSqlPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Synapse Sql Pool ID %q: %+v", input, err)
	}

	synapseSqlPool := SynapseSqlPoolId{
		Workspace: &SynapseWorkspaceId{
			SubscriptionID: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		},
	}
	if synapseSqlPool.Workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if synapseSqlPool.Name, err = id.PopSegment("sqlPools"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &synapseSqlPool, nil
}
