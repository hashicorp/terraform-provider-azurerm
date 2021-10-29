package webapplicationfirewallpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorWebApplicationFirewallPoliciesId struct {
	SubscriptionId                            string
	ResourceGroup                             string
	FrontDoorWebApplicationFirewallPolicyName string
}

func NewFrontDoorWebApplicationFirewallPoliciesID(subscriptionId, resourceGroup, frontDoorWebApplicationFirewallPolicyName string) FrontDoorWebApplicationFirewallPoliciesId {
	return FrontDoorWebApplicationFirewallPoliciesId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		FrontDoorWebApplicationFirewallPolicyName: frontDoorWebApplicationFirewallPolicyName,
	}
}

func (id FrontDoorWebApplicationFirewallPoliciesId) String() string {
	segments := []string{
		fmt.Sprintf("Front Door Web Application Firewall Policy Name %q", id.FrontDoorWebApplicationFirewallPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Web Application Firewall Policies", segmentsStr)
}

func (id FrontDoorWebApplicationFirewallPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FrontDoorWebApplicationFirewallPolicyName)
}

// ParseFrontDoorWebApplicationFirewallPoliciesID parses a FrontDoorWebApplicationFirewallPolicies ID into an FrontDoorWebApplicationFirewallPoliciesId struct
func ParseFrontDoorWebApplicationFirewallPoliciesID(input string) (*FrontDoorWebApplicationFirewallPoliciesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorWebApplicationFirewallPoliciesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment("frontDoorWebApplicationFirewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseFrontDoorWebApplicationFirewallPoliciesIDInsensitively parses an FrontDoorWebApplicationFirewallPolicies ID into an FrontDoorWebApplicationFirewallPoliciesId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseFrontDoorWebApplicationFirewallPoliciesID method should be used instead for validation etc.
func ParseFrontDoorWebApplicationFirewallPoliciesIDInsensitively(input string) (*FrontDoorWebApplicationFirewallPoliciesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorWebApplicationFirewallPoliciesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'frontDoorWebApplicationFirewallPolicies' segment
	frontDoorWebApplicationFirewallPoliciesKey := "frontDoorWebApplicationFirewallPolicies"
	for key := range id.Path {
		if strings.EqualFold(key, frontDoorWebApplicationFirewallPoliciesKey) {
			frontDoorWebApplicationFirewallPoliciesKey = key
			break
		}
	}
	if resourceId.FrontDoorWebApplicationFirewallPolicyName, err = id.PopSegment(frontDoorWebApplicationFirewallPoliciesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
