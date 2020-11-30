package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseWorkspaceId struct {
	SubscriptionID string
	ResourceGroup  string
	Name           string
}

func NewSynapseWorkspaceId(subscriptionId, resourceGroup, name string) SynapseWorkspaceId {
	return SynapseWorkspaceId{
		SubscriptionID: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func SynapseWorkspaceID(input string) (*SynapseWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapseWorkspace ID %q: %+v", input, err)
	}

	synapseWorkspace := SynapseWorkspaceId{
		ResourceGroup:  id.ResourceGroup,
		SubscriptionID: id.SubscriptionID,
	}
	if synapseWorkspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &synapseWorkspace, nil
}

func (id SynapseWorkspaceId) ID(_ string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s", id.SubscriptionID, id.ResourceGroup, id.Name)
}
