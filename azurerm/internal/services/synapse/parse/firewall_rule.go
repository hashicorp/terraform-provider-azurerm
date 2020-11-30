package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallRuleId struct {
	SubscriptionID string
	ResourceGroup  string
	WorkspaceName  string
	Name           string
}

func FirewallRuleID(input string) (*FirewallRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing synapseWorkspace ID %q: %+v", input, err)
	}

	FirewallRuleId := FirewallRuleId{
		SubscriptionID: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if FirewallRuleId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
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
