package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallPolicyPolicyId struct {
	ResourceGroup string
	Name          string
}

func FirewallPolicyPolicyID(input string) (*FirewallPolicyPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Firewall Policy Policy ID %q: %+v", input, err)
	}

	policy := FirewallPolicyPolicyId{
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
