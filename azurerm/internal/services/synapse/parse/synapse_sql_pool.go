package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseSqlPoolId struct {
	ResourceGroup string
	WorkspaceId   string
	WorkspaceName string
	Name          string
}

func SynapseSqlPoolID(input string) (*SynapseSqlPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Synapse Sql Pool ID %q: %+v", input, err)
	}

	synapseSqlPool := SynapseSqlPoolId{
		ResourceGroup: id.ResourceGroup,
	}
	if synapseSqlPool.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if synapseSqlPool.Name, err = id.PopSegment("sqlPools"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	synapseSqlPool.WorkspaceId = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s", id.SubscriptionID, id.ResourceGroup, synapseSqlPool.WorkspaceName)
	return &synapseSqlPool, nil
}
