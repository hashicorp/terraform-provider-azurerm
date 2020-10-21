package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceEndpointPolicyId struct {
	ResourceGroup string
	Name          string
}

func (id ServiceEndpointPolicyId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewServiceEndpointPolicyID(resourceGroup, name string) ServiceEndpointPolicyId {
	return ServiceEndpointPolicyId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func ServiceEndpointPolicyID(input string) (*ServiceEndpointPolicyId, error) {
	rawId, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Service Endpoint Policy ID %q: %+v", input, err)
	}

	id := ServiceEndpointPolicyId{
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

type ServiceEndpointPolicyDefinitionId struct {
	ResourceGroup string
	Policy        string
	Name          string
}

func (id ServiceEndpointPolicyDefinitionId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s/serviceEndpointPolicyDefinitions/%s",
		subscriptionId, id.ResourceGroup, id.Policy, id.Name)
}

func NewServiceEndpointPolicyDefinitionID(resourceGroup, policy, name string) ServiceEndpointPolicyDefinitionId {
	return ServiceEndpointPolicyDefinitionId{
		ResourceGroup: resourceGroup,
		Policy:        policy,
		Name:          name,
	}
}

func ServiceEndpointPolicyDefinitionID(input string) (*ServiceEndpointPolicyDefinitionId, error) {
	rawId, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Service Endpoint Policy Definition ID %q: %+v", input, err)
	}

	id := ServiceEndpointPolicyDefinitionId{
		ResourceGroup: rawId.ResourceGroup,
	}

	if id.Policy, err = rawId.PopSegment("serviceEndpointPolicies"); err != nil {
		return nil, err
	}

	if id.Name, err = rawId.PopSegment("serviceEndpointPolicyDefinitions"); err != nil {
		return nil, err
	}

	if err := rawId.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &id, nil
}
