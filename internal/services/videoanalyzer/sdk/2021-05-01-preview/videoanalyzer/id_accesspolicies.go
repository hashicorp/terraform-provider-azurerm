package videoanalyzer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AccessPoliciesId{}

// AccessPoliciesId is a struct representing the Resource ID for a Access Policies
type AccessPoliciesId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	AccessPolicyName  string
}

// NewAccessPoliciesID returns a new AccessPoliciesId struct
func NewAccessPoliciesID(subscriptionId string, resourceGroupName string, accountName string, accessPolicyName string) AccessPoliciesId {
	return AccessPoliciesId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		AccessPolicyName:  accessPolicyName,
	}
}

// ParseAccessPoliciesID parses 'input' into a AccessPoliciesId
func ParseAccessPoliciesID(input string) (*AccessPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessPoliciesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.AccessPolicyName, ok = parsed.Parsed["accessPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'accessPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAccessPoliciesIDInsensitively parses 'input' case-insensitively into a AccessPoliciesId
// note: this method should only be used for API response data and not user input
func ParseAccessPoliciesIDInsensitively(input string) (*AccessPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessPoliciesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.AccessPolicyName, ok = parsed.Parsed["accessPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'accessPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAccessPoliciesID checks that 'input' can be parsed as a Access Policies ID
func ValidateAccessPoliciesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccessPoliciesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Access Policies ID
func (id AccessPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s/accessPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.AccessPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Access Policies ID
func (id AccessPoliciesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticVideoAnalyzers", "videoAnalyzers", "videoAnalyzers"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticAccessPolicies", "accessPolicies", "accessPolicies"),
		resourceids.UserSpecifiedSegment("accessPolicyName", "accessPolicyValue"),
	}
}

// String returns a human-readable description of this Access Policies ID
func (id AccessPoliciesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Access Policy Name: %q", id.AccessPolicyName),
	}
	return fmt.Sprintf("Access Policies (%s)", strings.Join(components, "\n"))
}
