package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type WebApplicationFirewallPolicyId struct {
	ResourceGroup string
	Name          string
}

func NewWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name string) WebApplicationFirewallPolicyId {
	return WebApplicationFirewallPolicyId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id WebApplicationFirewallPolicyId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func WebApplicationFirewallPolicyID(input string) (*WebApplicationFirewallPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Web Application Firewall Policy ID %q: %+v", input, err)
	}

	policy := WebApplicationFirewallPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if policy.Name, err = id.PopSegment("frontDoorWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &policy, nil
}
