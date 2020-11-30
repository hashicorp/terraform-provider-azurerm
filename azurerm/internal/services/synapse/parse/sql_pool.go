package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SqlPoolId struct {
	Workspace *WorkspaceId
	Name      string
}

func SqlPoolID(input string) (*SqlPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Synapse Sql Pool ID %q: %+v", input, err)
	}

	synapseSqlPool := SqlPoolId{
		Workspace: &WorkspaceId{
			SubscriptionId: id.SubscriptionID,
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
