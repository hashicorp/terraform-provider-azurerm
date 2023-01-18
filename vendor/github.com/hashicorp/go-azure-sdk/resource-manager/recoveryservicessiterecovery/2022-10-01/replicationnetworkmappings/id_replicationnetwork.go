package replicationnetworkmappings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ReplicationNetworkId{}

// ReplicationNetworkId is a struct representing the Resource ID for a Replication Network
type ReplicationNetworkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
	FabricName        string
	NetworkName       string
}

// NewReplicationNetworkID returns a new ReplicationNetworkId struct
func NewReplicationNetworkID(subscriptionId string, resourceGroupName string, resourceName string, fabricName string, networkName string) ReplicationNetworkId {
	return ReplicationNetworkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceName:      resourceName,
		FabricName:        fabricName,
		NetworkName:       networkName,
	}
}

// ParseReplicationNetworkID parses 'input' into a ReplicationNetworkId
func ParseReplicationNetworkID(input string) (*ReplicationNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationNetworkId{}

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

	if id.NetworkName, ok = parsed.Parsed["networkName"]; !ok {
		return nil, fmt.Errorf("the segment 'networkName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseReplicationNetworkIDInsensitively parses 'input' case-insensitively into a ReplicationNetworkId
// note: this method should only be used for API response data and not user input
func ParseReplicationNetworkIDInsensitively(input string) (*ReplicationNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationNetworkId{}

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

	if id.NetworkName, ok = parsed.Parsed["networkName"]; !ok {
		return nil, fmt.Errorf("the segment 'networkName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateReplicationNetworkID checks that 'input' can be parsed as a Replication Network ID
func ValidateReplicationNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Network ID
func (id ReplicationNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.FabricName, id.NetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Network ID
func (id ReplicationNetworkId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticReplicationNetworks", "replicationNetworks", "replicationNetworks"),
		resourceids.UserSpecifiedSegment("networkName", "networkValue"),
	}
}

// String returns a human-readable description of this Replication Network ID
func (id ReplicationNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Fabric Name: %q", id.FabricName),
		fmt.Sprintf("Network Name: %q", id.NetworkName),
	}
	return fmt.Sprintf("Replication Network (%s)", strings.Join(components, "\n"))
}
