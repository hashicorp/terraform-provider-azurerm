package securitypolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SecurityPoliciesId{}

// SecurityPoliciesId is a struct representing the Resource ID for a Security Policies
type SecurityPoliciesId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ProfileName        string
	SecurityPolicyName string
}

// NewSecurityPoliciesID returns a new SecurityPoliciesId struct
func NewSecurityPoliciesID(subscriptionId string, resourceGroupName string, profileName string, securityPolicyName string) SecurityPoliciesId {
	return SecurityPoliciesId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ProfileName:        profileName,
		SecurityPolicyName: securityPolicyName,
	}
}

// ParseSecurityPoliciesID parses 'input' into a SecurityPoliciesId
func ParseSecurityPoliciesID(input string) (*SecurityPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecurityPoliciesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecurityPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.SecurityPolicyName, ok = parsed.Parsed["securityPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'securityPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSecurityPoliciesIDInsensitively parses 'input' case-insensitively into a SecurityPoliciesId
// note: this method should only be used for API response data and not user input
func ParseSecurityPoliciesIDInsensitively(input string) (*SecurityPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(SecurityPoliciesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SecurityPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProfileName, ok = parsed.Parsed["profileName"]; !ok {
		return nil, fmt.Errorf("the segment 'profileName' was not found in the resource id %q", input)
	}

	if id.SecurityPolicyName, ok = parsed.Parsed["securityPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'securityPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSecurityPoliciesID checks that 'input' can be parsed as a Security Policies ID
func ValidateSecurityPoliciesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSecurityPoliciesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Security Policies ID
func (id SecurityPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CDN/profiles/%s/securityPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.SecurityPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Security Policies ID
func (id SecurityPoliciesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCDN", "Microsoft.CDN", "Microsoft.CDN"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileValue"),
		resourceids.StaticSegment("staticSecurityPolicies", "securityPolicies", "securityPolicies"),
		resourceids.UserSpecifiedSegment("securityPolicyName", "securityPolicyValue"),
	}
}

// String returns a human-readable description of this Security Policies ID
func (id SecurityPoliciesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Security Policy Name: %q", id.SecurityPolicyName),
	}
	return fmt.Sprintf("Security Policies (%s)", strings.Join(components, "\n"))
}
