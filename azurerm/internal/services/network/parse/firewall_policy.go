package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallPolicyId struct {
	ResourceGroup string
	Name          string
}

func (id FirewallPolicyId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewFirewallPolicyID(resourceGroup, name string) FirewallPolicyId {
	return FirewallPolicyId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func FirewallPolicyID(input string) (*FirewallPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Firewall Policy ID %q: %+v", input, err)
	}

	policy := FirewallPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if policy.Name, err = id.PopSegment("firewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &policy, nil
}
