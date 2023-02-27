package replicationprotectableitems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ReplicationProtectableItemId{}

// ReplicationProtectableItemId is a struct representing the Resource ID for a Replication Protectable Item
type ReplicationProtectableItemId struct {
	SubscriptionId                     string
	ResourceGroupName                  string
	VaultName                          string
	ReplicationFabricName              string
	ReplicationProtectionContainerName string
	ReplicationProtectableItemName     string
}

// NewReplicationProtectableItemID returns a new ReplicationProtectableItemId struct
func NewReplicationProtectableItemID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationProtectionContainerName string, replicationProtectableItemName string) ReplicationProtectableItemId {
	return ReplicationProtectableItemId{
		SubscriptionId:                     subscriptionId,
		ResourceGroupName:                  resourceGroupName,
		VaultName:                          vaultName,
		ReplicationFabricName:              replicationFabricName,
		ReplicationProtectionContainerName: replicationProtectionContainerName,
		ReplicationProtectableItemName:     replicationProtectableItemName,
	}
}

// ParseReplicationProtectableItemID parses 'input' into a ReplicationProtectableItemId
func ParseReplicationProtectableItemID(input string) (*ReplicationProtectableItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectableItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectableItemId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	if id.ReplicationFabricName, ok = parsed.Parsed["replicationFabricName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicationFabricName' was not found in the resource id %q", input)
	}

	if id.ReplicationProtectionContainerName, ok = parsed.Parsed["replicationProtectionContainerName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicationProtectionContainerName' was not found in the resource id %q", input)
	}

	if id.ReplicationProtectableItemName, ok = parsed.Parsed["replicationProtectableItemName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicationProtectableItemName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseReplicationProtectableItemIDInsensitively parses 'input' case-insensitively into a ReplicationProtectableItemId
// note: this method should only be used for API response data and not user input
func ParseReplicationProtectableItemIDInsensitively(input string) (*ReplicationProtectableItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectableItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectableItemId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	if id.ReplicationFabricName, ok = parsed.Parsed["replicationFabricName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicationFabricName' was not found in the resource id %q", input)
	}

	if id.ReplicationProtectionContainerName, ok = parsed.Parsed["replicationProtectionContainerName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicationProtectionContainerName' was not found in the resource id %q", input)
	}

	if id.ReplicationProtectableItemName, ok = parsed.Parsed["replicationProtectableItemName"]; !ok {
		return nil, fmt.Errorf("the segment 'replicationProtectableItemName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateReplicationProtectableItemID checks that 'input' can be parsed as a Replication Protectable Item ID
func ValidateReplicationProtectableItemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationProtectableItemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Protectable Item ID
func (id ReplicationProtectableItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s/replicationProtectableItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectableItemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Protectable Item ID
func (id ReplicationProtectableItemId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricValue"),
		resourceids.StaticSegment("staticReplicationProtectionContainers", "replicationProtectionContainers", "replicationProtectionContainers"),
		resourceids.UserSpecifiedSegment("replicationProtectionContainerName", "replicationProtectionContainerValue"),
		resourceids.StaticSegment("staticReplicationProtectableItems", "replicationProtectableItems", "replicationProtectableItems"),
		resourceids.UserSpecifiedSegment("replicationProtectableItemName", "replicationProtectableItemValue"),
	}
}

// String returns a human-readable description of this Replication Protectable Item ID
func (id ReplicationProtectableItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Protection Container Name: %q", id.ReplicationProtectionContainerName),
		fmt.Sprintf("Replication Protectable Item Name: %q", id.ReplicationProtectableItemName),
	}
	return fmt.Sprintf("Replication Protectable Item (%s)", strings.Join(components, "\n"))
}
