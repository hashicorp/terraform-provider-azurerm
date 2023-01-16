package replicationprotectioncontainermappings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ReplicationProtectionContainerMappingId{}

// ReplicationProtectionContainerMappingId is a struct representing the Resource ID for a Replication Protection Container Mapping
type ReplicationProtectionContainerMappingId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ResourceName            string
	FabricName              string
	ProtectionContainerName string
	MappingName             string
}

// NewReplicationProtectionContainerMappingID returns a new ReplicationProtectionContainerMappingId struct
func NewReplicationProtectionContainerMappingID(subscriptionId string, resourceGroupName string, resourceName string, fabricName string, protectionContainerName string, mappingName string) ReplicationProtectionContainerMappingId {
	return ReplicationProtectionContainerMappingId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ResourceName:            resourceName,
		FabricName:              fabricName,
		ProtectionContainerName: protectionContainerName,
		MappingName:             mappingName,
	}
}

// ParseReplicationProtectionContainerMappingID parses 'input' into a ReplicationProtectionContainerMappingId
func ParseReplicationProtectionContainerMappingID(input string) (*ReplicationProtectionContainerMappingId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectionContainerMappingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectionContainerMappingId{}

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

	if id.MappingName, ok = parsed.Parsed["mappingName"]; !ok {
		return nil, fmt.Errorf("the segment 'mappingName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseReplicationProtectionContainerMappingIDInsensitively parses 'input' case-insensitively into a ReplicationProtectionContainerMappingId
// note: this method should only be used for API response data and not user input
func ParseReplicationProtectionContainerMappingIDInsensitively(input string) (*ReplicationProtectionContainerMappingId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectionContainerMappingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectionContainerMappingId{}

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

	if id.MappingName, ok = parsed.Parsed["mappingName"]; !ok {
		return nil, fmt.Errorf("the segment 'mappingName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateReplicationProtectionContainerMappingID checks that 'input' can be parsed as a Replication Protection Container Mapping ID
func ValidateReplicationProtectionContainerMappingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationProtectionContainerMappingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Protection Container Mapping ID
func (id ReplicationProtectionContainerMappingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s/replicationProtectionContainerMappings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.FabricName, id.ProtectionContainerName, id.MappingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Protection Container Mapping ID
func (id ReplicationProtectionContainerMappingId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticReplicationProtectionContainerMappings", "replicationProtectionContainerMappings", "replicationProtectionContainerMappings"),
		resourceids.UserSpecifiedSegment("mappingName", "mappingValue"),
	}
}

// String returns a human-readable description of this Replication Protection Container Mapping ID
func (id ReplicationProtectionContainerMappingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Fabric Name: %q", id.FabricName),
		fmt.Sprintf("Protection Container Name: %q", id.ProtectionContainerName),
		fmt.Sprintf("Mapping Name: %q", id.MappingName),
	}
	return fmt.Sprintf("Replication Protection Container Mapping (%s)", strings.Join(components, "\n"))
}
