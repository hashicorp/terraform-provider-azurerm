package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubnetServiceEndpointPolicyId struct {
	SubscriptionId            string
	ResourceGroup             string
	ServiceEndpointPolicyName string
}

func NewSubnetServiceEndpointPolicyID(subscriptionId, resourceGroup, serviceEndpointPolicyName string) SubnetServiceEndpointPolicyId {
	return SubnetServiceEndpointPolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		ServiceEndpointPolicyName: serviceEndpointPolicyName,
	}
}

func (id SubnetServiceEndpointPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Service Endpoint Policy Name %q", id.ServiceEndpointPolicyName),
	}
	return strings.Join(segments, " / ")
}

func (id SubnetServiceEndpointPolicyId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/serviceEndpointPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceEndpointPolicyName)
}

// SubnetServiceEndpointPolicyID parses a SubnetServiceEndpointPolicy ID into an SubnetServiceEndpointPolicyId struct
func SubnetServiceEndpointPolicyID(input string) (*SubnetServiceEndpointPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubnetServiceEndpointPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceEndpointPolicyName, err = id.PopSegment("serviceEndpointPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
