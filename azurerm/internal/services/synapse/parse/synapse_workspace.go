package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseWorkspaceId struct {
	ResourceGroup string
	Name          string
}

func SynapseWorkspaceID(input string) (*SynapseWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapseWorkspace ID %q: %+v", input, err)
	}

	synapseWorkspace := SynapseWorkspaceId{
		ResourceGroup: id.ResourceGroup,
	}
	if synapseWorkspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &synapseWorkspace, nil
}
