package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseBigDataPoolId struct {
	ResourceGroup string
	WorkspaceName string
	WorkspaceId   string
	Name          string
}

func SynapseBigDataPoolID(input string) (*SynapseBigDataPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapseBigDataPool ID %q: %+v", input, err)
	}

	synapseBigDataPool := SynapseBigDataPoolId{
		ResourceGroup: id.ResourceGroup,
	}
	if synapseBigDataPool.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if synapseBigDataPool.Name, err = id.PopSegment("bigDataPools"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	synapseBigDataPool.WorkspaceId = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s", id.SubscriptionID, id.ResourceGroup, synapseBigDataPool.WorkspaceName)
	return &synapseBigDataPool, nil
}
