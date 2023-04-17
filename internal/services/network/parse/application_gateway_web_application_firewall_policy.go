package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ApplicationGatewayWebApplicationFirewallPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewApplicationGatewayWebApplicationFirewallPolicyID(subscriptionId, resourceGroup, name string) ApplicationGatewayWebApplicationFirewallPolicyId {
	return ApplicationGatewayWebApplicationFirewallPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id ApplicationGatewayWebApplicationFirewallPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Application Gateway Web Application Firewall Policy", segmentsStr)
}

func (id ApplicationGatewayWebApplicationFirewallPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// ApplicationGatewayWebApplicationFirewallPolicyID parses a ApplicationGatewayWebApplicationFirewallPolicy ID into an ApplicationGatewayWebApplicationFirewallPolicyId struct
func ApplicationGatewayWebApplicationFirewallPolicyID(input string) (*ApplicationGatewayWebApplicationFirewallPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayWebApplicationFirewallPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("ApplicationGatewayWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ApplicationGatewayWebApplicationFirewallPolicyIDInsensitively parses an ApplicationGatewayWebApplicationFirewallPolicy ID into an ApplicationGatewayWebApplicationFirewallPolicyId struct, insensitively
// This should only be used to parse an ID for rewriting, the ApplicationGatewayWebApplicationFirewallPolicyID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func ApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(input string) (*ApplicationGatewayWebApplicationFirewallPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayWebApplicationFirewallPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'ApplicationGatewayWebApplicationFirewallPolicies' segment
	ApplicationGatewayWebApplicationFirewallPoliciesKey := "ApplicationGatewayWebApplicationFirewallPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, ApplicationGatewayWebApplicationFirewallPoliciesKey) {
			ApplicationGatewayWebApplicationFirewallPoliciesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(ApplicationGatewayWebApplicationFirewallPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
