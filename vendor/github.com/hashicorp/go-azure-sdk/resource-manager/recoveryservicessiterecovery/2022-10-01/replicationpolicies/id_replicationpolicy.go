package replicationpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ReplicationPolicyId{}

// ReplicationPolicyId is a struct representing the Resource ID for a Replication Policy
type ReplicationPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
	PolicyName        string
}

// NewReplicationPolicyID returns a new ReplicationPolicyId struct
func NewReplicationPolicyID(subscriptionId string, resourceGroupName string, resourceName string, policyName string) ReplicationPolicyId {
	return ReplicationPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceName:      resourceName,
		PolicyName:        policyName,
	}
}

// ParseReplicationPolicyID parses 'input' into a ReplicationPolicyId
func ParseReplicationPolicyID(input string) (*ReplicationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.PolicyName, ok = parsed.Parsed["policyName"]; !ok {
		return nil, fmt.Errorf("the segment 'policyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseReplicationPolicyIDInsensitively parses 'input' case-insensitively into a ReplicationPolicyId
// note: this method should only be used for API response data and not user input
func ParseReplicationPolicyIDInsensitively(input string) (*ReplicationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.PolicyName, ok = parsed.Parsed["policyName"]; !ok {
		return nil, fmt.Errorf("the segment 'policyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateReplicationPolicyID checks that 'input' can be parsed as a Replication Policy ID
func ValidateReplicationPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Policy ID
func (id ReplicationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.PolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Policy ID
func (id ReplicationPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticReplicationPolicies", "replicationPolicies", "replicationPolicies"),
		resourceids.UserSpecifiedSegment("policyName", "policyValue"),
	}
}

// String returns a human-readable description of this Replication Policy ID
func (id ReplicationPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Policy Name: %q", id.PolicyName),
	}
	return fmt.Sprintf("Replication Policy (%s)", strings.Join(components, "\n"))
}
