package replicationprotecteditems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ReplicationProtectedItemId{}

// ReplicationProtectedItemId is a struct representing the Resource ID for a Replication Protected Item
type ReplicationProtectedItemId struct {
	SubscriptionId              string
	ResourceGroupName           string
	ResourceName                string
	FabricName                  string
	ProtectionContainerName     string
	ReplicatedProtectedItemName string
}

// NewReplicationProtectedItemID returns a new ReplicationProtectedItemId struct
func NewReplicationProtectedItemID(subscriptionId string, resourceGroupName string, resourceName string, fabricName string, protectionContainerName string, replicatedProtectedItemName string) ReplicationProtectedItemId {
	return ReplicationProtectedItemId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		ResourceName:                resourceName,
		FabricName:                  fabricName,
		ProtectionContainerName:     protectionContainerName,
		ReplicatedProtectedItemName: replicatedProtectedItemName,
	}
}

// ParseReplicationProtectedItemID parses 'input' into a ReplicationProtectedItemId
func ParseReplicationProtectedItemID(input string) (*ReplicationProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectedItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectedItemId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.FabricName, ok = parsed.Parsed["fabricName"]; !ok {
		return nil, fmt.Errorf("the segment 'fabricName' was not found in the resource id %q", input)
	}

	if id.ProtectionContainerName, ok = parsed.Parsed["protectionContainerName"]; !ok {
		return nil, fmt.Errorf("the segment 'protectionContainerName' was not found in the resource id %q", input)
	}

	if id.ReplicatedProtectedItemName, ok = parsed.Parsed["replicatedProtectedItemName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicatedProtectedItemName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseReplicationProtectedItemIDInsensitively parses 'input' case-insensitively into a ReplicationProtectedItemId
// note: this method should only be used for API response data and not user input
func ParseReplicationProtectedItemIDInsensitively(input string) (*ReplicationProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectedItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectedItemId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.FabricName, ok = parsed.Parsed["fabricName"]; !ok {
		return nil, fmt.Errorf("the segment 'fabricName' was not found in the resource id %q", input)
	}

	if id.ProtectionContainerName, ok = parsed.Parsed["protectionContainerName"]; !ok {
		return nil, fmt.Errorf("the segment 'protectionContainerName' was not found in the resource id %q", input)
	}

	if id.ReplicatedProtectedItemName, ok = parsed.Parsed["replicatedProtectedItemName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicatedProtectedItemName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateReplicationProtectedItemID checks that 'input' can be parsed as a Replication Protected Item ID
func ValidateReplicationProtectedItemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationProtectedItemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Protected Item ID
func (id ReplicationProtectedItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s/replicationProtectedItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.FabricName, id.ProtectionContainerName, id.ReplicatedProtectedItemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Protected Item ID
func (id ReplicationProtectedItemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("fabricName", "fabricValue"),
		resourceids.StaticSegment("staticReplicationProtectionContainers", "replicationProtectionContainers", "replicationProtectionContainers"),
		resourceids.UserSpecifiedSegment("protectionContainerName", "protectionContainerValue"),
		resourceids.StaticSegment("staticReplicationProtectedItems", "replicationProtectedItems", "replicationProtectedItems"),
		resourceids.UserSpecifiedSegment("replicatedProtectedItemName", "replicatedProtectedItemValue"),
	}
}

// String returns a human-readable description of this Replication Protected Item ID
func (id ReplicationProtectedItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Fabric Name: %q", id.FabricName),
		fmt.Sprintf("Protection Container Name: %q", id.ProtectionContainerName),
		fmt.Sprintf("Replicated Protected Item Name: %q", id.ReplicatedProtectedItemName),
	}
	return fmt.Sprintf("Replication Protected Item (%s)", strings.Join(components, "\n"))
}
