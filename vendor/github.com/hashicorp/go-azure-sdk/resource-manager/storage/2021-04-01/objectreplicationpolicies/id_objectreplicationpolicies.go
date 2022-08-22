package objectreplicationpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ObjectReplicationPoliciesId{}

// ObjectReplicationPoliciesId is a struct representing the Resource ID for a Object Replication Policies
type ObjectReplicationPoliciesId struct {
	SubscriptionId            string
	ResourceGroupName         string
	AccountName               string
	ObjectReplicationPolicyId string
}

// NewObjectReplicationPoliciesID returns a new ObjectReplicationPoliciesId struct
func NewObjectReplicationPoliciesID(subscriptionId string, resourceGroupName string, accountName string, objectReplicationPolicyId string) ObjectReplicationPoliciesId {
	return ObjectReplicationPoliciesId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		AccountName:               accountName,
		ObjectReplicationPolicyId: objectReplicationPolicyId,
	}
}

// ParseObjectReplicationPoliciesID parses 'input' into a ObjectReplicationPoliciesId
func ParseObjectReplicationPoliciesID(input string) (*ObjectReplicationPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(ObjectReplicationPoliciesId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ObjectReplicationPoliciesId{}

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

// ParseObjectReplicationPoliciesIDInsensitively parses 'input' case-insensitively into a ObjectReplicationPoliciesId
// note: this method should only be used for API response data and not user input
func ParseObjectReplicationPoliciesIDInsensitively(input string) (*ObjectReplicationPoliciesId, error) {
	parser := resourceids.NewParserFromResourceIdType(ObjectReplicationPoliciesId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ObjectReplicationPoliciesId{}

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

// ValidateObjectReplicationPoliciesID checks that 'input' can be parsed as a Object Replication Policies ID
func ValidateObjectReplicationPoliciesID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseObjectReplicationPoliciesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Object Replication Policies ID
func (id ObjectReplicationPoliciesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/objectReplicationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ObjectReplicationPolicyId)
}

// Segments returns a slice of Resource ID Segments which comprise this Object Replication Policies ID
func (id ObjectReplicationPoliciesId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Object Replication Policies ID
func (id ObjectReplicationPoliciesId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Object Replication Policy: %q", id.ObjectReplicationPolicyId),
	}
	return fmt.Sprintf("Object Replication Policies (%s)", strings.Join(components, "\n"))
}
