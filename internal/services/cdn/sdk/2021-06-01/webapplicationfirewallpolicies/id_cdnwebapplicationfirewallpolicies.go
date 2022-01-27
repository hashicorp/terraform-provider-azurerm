package webapplicationfirewallpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CdnWebApplicationFirewallPoliciesId{}

// CdnWebApplicationFirewallPoliciesId is a struct representing the Resource ID for a Cdn Web Application Firewall Policies
type CdnWebApplicationFirewallPoliciesId struct {
	SubscriptionId    string
	ResourceGroupName string
	PolicyName        string
}

// NewCdnWebApplicationFirewallPoliciesID returns a new CdnWebApplicationFirewallPoliciesId struct
func NewCdnWebApplicationFirewallPoliciesID(subscriptionId string, resourceGroupName string, policyName string) CdnWebApplicationFirewallPoliciesId {
	return CdnWebApplicationFirewallPoliciesId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PolicyName:        policyName,
	}
}

// ParseCdnWebApplicationFirewallPoliciesID parses 'input' into a CdnWebApplicationFirewallPoliciesId
func ParseCdnWebApplicationFirewallPoliciesID(input string) (*CdnWebApplicationFirewallPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(CdnWebApplicationFirewallPoliciesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CdnWebApplicationFirewallPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PolicyName, ok = parsed.Parsed["policyName"]; !ok {
		return nil, fmt.Errorf("the segment 'policyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCdnWebApplicationFirewallPoliciesIDInsensitively parses 'input' case-insensitively into a CdnWebApplicationFirewallPoliciesId
// note: this method should only be used for API response data and not user input
func ParseCdnWebApplicationFirewallPoliciesIDInsensitively(input string) (*CdnWebApplicationFirewallPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(CdnWebApplicationFirewallPoliciesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CdnWebApplicationFirewallPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PolicyName, ok = parsed.Parsed["policyName"]; !ok {
		return nil, fmt.Errorf("the segment 'policyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCdnWebApplicationFirewallPoliciesID checks that 'input' can be parsed as a Cdn Web Application Firewall Policies ID
func ValidateCdnWebApplicationFirewallPoliciesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCdnWebApplicationFirewallPoliciesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cdn Web Application Firewall Policies ID
func (id CdnWebApplicationFirewallPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/cdnWebApplicationFirewallPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cdn Web Application Firewall Policies ID
func (id CdnWebApplicationFirewallPoliciesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticCdnWebApplicationFirewallPolicies", "cdnWebApplicationFirewallPolicies", "cdnWebApplicationFirewallPolicies"),
		resourceids.UserSpecifiedSegment("policyName", "policyValue"),
	}
}

// String returns a human-readable description of this Cdn Web Application Firewall Policies ID
func (id CdnWebApplicationFirewallPoliciesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Policy Name: %q", id.PolicyName),
	}
	return fmt.Sprintf("Cdn Web Application Firewall Policies (%s)", strings.Join(components, "\n"))
}
