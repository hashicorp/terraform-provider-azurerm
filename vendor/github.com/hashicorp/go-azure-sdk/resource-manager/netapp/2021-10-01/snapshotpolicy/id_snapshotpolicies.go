package snapshotpolicy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SnapshotPoliciesId{}

// SnapshotPoliciesId is a struct representing the Resource ID for a Snapshot Policies
type SnapshotPoliciesId struct {
	SubscriptionId     string
	ResourceGroupName  string
	AccountName        string
	SnapshotPolicyName string
}

// NewSnapshotPoliciesID returns a new SnapshotPoliciesId struct
func NewSnapshotPoliciesID(subscriptionId string, resourceGroupName string, accountName string, snapshotPolicyName string) SnapshotPoliciesId {
	return SnapshotPoliciesId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		AccountName:        accountName,
		SnapshotPolicyName: snapshotPolicyName,
	}
}

// ParseSnapshotPoliciesID parses 'input' into a SnapshotPoliciesId
func ParseSnapshotPoliciesID(input string) (*SnapshotPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(SnapshotPoliciesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SnapshotPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.SnapshotPolicyName, ok = parsed.Parsed["snapshotPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'snapshotPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSnapshotPoliciesIDInsensitively parses 'input' case-insensitively into a SnapshotPoliciesId
// note: this method should only be used for API response data and not user input
func ParseSnapshotPoliciesIDInsensitively(input string) (*SnapshotPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(SnapshotPoliciesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SnapshotPoliciesId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.SnapshotPolicyName, ok = parsed.Parsed["snapshotPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'snapshotPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSnapshotPoliciesID checks that 'input' can be parsed as a Snapshot Policies ID
func ValidateSnapshotPoliciesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSnapshotPoliciesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Snapshot Policies ID
func (id SnapshotPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/snapshotPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.SnapshotPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Snapshot Policies ID
func (id SnapshotPoliciesId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticSnapshotPolicies", "snapshotPolicies", "snapshotPolicies"),
		resourceids.UserSpecifiedSegment("snapshotPolicyName", "snapshotPolicyValue"),
	}
}

// String returns a human-readable description of this Snapshot Policies ID
func (id SnapshotPoliciesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Snapshot Policy Name: %q", id.SnapshotPolicyName),
	}
	return fmt.Sprintf("Snapshot Policies (%s)", strings.Join(components, "\n"))
}
