package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type JitNetworkAccessPolicyId struct {
	ResourceGroup string
	Name          string
	Location      string
}

func NewJitNetworkAccessPolicyId(resourceGroup, name, location string) JitNetworkAccessPolicyId {
	return JitNetworkAccessPolicyId{
		ResourceGroup: resourceGroup,
		Name:          name,
		Location:      location,
	}
}

func (id JitNetworkAccessPolicyId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Security/locations/%s/jitNetworkAccessPolicies/%s", subscriptionId, id.ResourceGroup, id.Location, id.Name)
}

func JitNetworkAccessPolicyID(input string) (*JitNetworkAccessPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Jit Network Access Policy ID %q: %+v", input, err)
	}

	jnp := JitNetworkAccessPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if jnp.Name, err = id.PopSegment("jitNetworkAccessPolicies"); err != nil {
		return nil, err
	}

	if jnp.Location, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &jnp, nil
}
