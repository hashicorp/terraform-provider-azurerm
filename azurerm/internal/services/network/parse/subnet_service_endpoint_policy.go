package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubnetServiceEndpointPolicyId struct {
	ResourceGroup string
	Name          string
}

func (id SubnetServiceEndpointPolicyId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewServiceEndpointPolicyID(resourceGroup, name string) SubnetServiceEndpointPolicyId {
	return SubnetServiceEndpointPolicyId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func SubnetServiceEndpointPolicyID(input string) (*SubnetServiceEndpointPolicyId, error) {
	rawId, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Subnet Service Endpoint Policy ID %q: %+v", input, err)
	}

	id := SubnetServiceEndpointPolicyId{
		ResourceGroup: rawId.ResourceGroup,
	}

	if id.Name, err = rawId.PopSegment("serviceEndpointPolicies"); err != nil {
		return nil, err
	}

	if err := rawId.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &id, nil
}
