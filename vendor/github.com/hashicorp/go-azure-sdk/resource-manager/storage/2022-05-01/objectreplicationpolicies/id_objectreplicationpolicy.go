package objectreplicationpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ObjectReplicationPolicyId{}

// ObjectReplicationPolicyId is a struct representing the Resource ID for a Object Replication Policy
type ObjectReplicationPolicyId struct {
	SubscriptionId            string
	ResourceGroupName         string
	AccountName               string
	ObjectReplicationPolicyId string
}

// NewObjectReplicationPolicyID returns a new ObjectReplicationPolicyId struct
func NewObjectReplicationPolicyID(subscriptionId string, resourceGroupName string, accountName string, objectReplicationPolicyId string) ObjectReplicationPolicyId {
	return ObjectReplicationPolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		AccountName:               accountName,
		ObjectReplicationPolicyId: objectReplicationPolicyId,
	}
}

// ParseObjectReplicationPolicyID parses 'input' into a ObjectReplicationPolicyId
func ParseObjectReplicationPolicyID(input string) (*ObjectReplicationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ObjectReplicationPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ObjectReplicationPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.ObjectReplicationPolicyId, ok = parsed.Parsed["objectReplicationPolicyId"]; !ok {
		return nil, fmt.Errorf("the segment 'objectReplicationPolicyId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseObjectReplicationPolicyIDInsensitively parses 'input' case-insensitively into a ObjectReplicationPolicyId
// note: this method should only be used for API response data and not user input
func ParseObjectReplicationPolicyIDInsensitively(input string) (*ObjectReplicationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ObjectReplicationPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ObjectReplicationPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.ObjectReplicationPolicyId, ok = parsed.Parsed["objectReplicationPolicyId"]; !ok {
		return nil, fmt.Errorf("the segment 'objectReplicationPolicyId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateObjectReplicationPolicyID checks that 'input' can be parsed as a Object Replication Policy ID
func ValidateObjectReplicationPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseObjectReplicationPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Object Replication Policy ID
func (id ObjectReplicationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ObjectReplicationPolicyId)
}

// Segments returns a slice of Resource ID Segments which comprise this Object Replication Policy ID
func (id ObjectReplicationPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticObjectReplicationPolicies", "objectReplicationPolicies", "objectReplicationPolicies"),
		resourceids.UserSpecifiedSegment("objectReplicationPolicyId", "objectReplicationPolicyIdValue"),
	}
}

// String returns a human-readable description of this Object Replication Policy ID
func (id ObjectReplicationPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Object Replication Policy: %q", id.ObjectReplicationPolicyId),
	}
	return fmt.Sprintf("Object Replication Policy (%s)", strings.Join(components, "\n"))
}
