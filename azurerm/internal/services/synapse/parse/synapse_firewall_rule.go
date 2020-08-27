package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SynapseFirewallRuleId struct {
	Workspace *SynapseWorkspaceId
	Name      string
}

func SynapseFirewallRuleID(input string) (*SynapseFirewallRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapseWorkspace ID %q: %+v", input, err)
	}

	FirewallRuleId := SynapseFirewallRuleId{
		Workspace: &SynapseWorkspaceId{
			SubscriptionID: id.SubscriptionID,
			ResourceGroup:  id.ResourceGroup,
		},
	}
	if FirewallRuleId.Workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if FirewallRuleId.Name, err = id.PopSegment("firewallRules"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &FirewallRuleId, nil
}
