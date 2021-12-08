package computepolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ComputePoliciesId{}

// ComputePoliciesId is a struct representing the Resource ID for a Compute Policies
type ComputePoliciesId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	ComputePolicyName string
}

// NewComputePoliciesID returns a new ComputePoliciesId struct
func NewComputePoliciesID(subscriptionId string, resourceGroupName string, accountName string, computePolicyName string) ComputePoliciesId {
	return ComputePoliciesId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		ComputePolicyName: computePolicyName,
	}
}

// ParseComputePoliciesID parses 'input' into a ComputePoliciesId
func ParseComputePoliciesID(input string) (*ComputePoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(ComputePoliciesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ComputePoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.ComputePolicyName, ok = parsed.Parsed["computePolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'computePolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseComputePoliciesIDInsensitively parses 'input' case-insensitively into a ComputePoliciesId
// note: this method should only be used for API response data and not user input
func ParseComputePoliciesIDInsensitively(input string) (*ComputePoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(ComputePoliciesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ComputePoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.ComputePolicyName, ok = parsed.Parsed["computePolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'computePolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateComputePoliciesID checks that 'input' can be parsed as a Compute Policies ID
func ValidateComputePoliciesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseComputePoliciesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Compute Policies ID
func (id ComputePoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeAnalytics/accounts/%s/computePolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ComputePolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Compute Policies ID
func (id ComputePoliciesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataLakeAnalytics", "Microsoft.DataLakeAnalytics", "Microsoft.DataLakeAnalytics"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticComputePolicies", "computePolicies", "computePolicies"),
		resourceids.UserSpecifiedSegment("computePolicyName", "computePolicyValue"),
	}
}

// String returns a human-readable description of this Compute Policies ID
func (id ComputePoliciesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Compute Policy Name: %q", id.ComputePolicyName),
	}
	return fmt.Sprintf("Compute Policies (%s)", strings.Join(components, "\n"))
}
