package replicationnetworks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationNetworkId{}

// ReplicationNetworkId is a struct representing the Resource ID for a Replication Network
type ReplicationNetworkId struct {
	SubscriptionId         string
	ResourceGroupName      string
	VaultName              string
	ReplicationFabricName  string
	ReplicationNetworkName string
}

// NewReplicationNetworkID returns a new ReplicationNetworkId struct
func NewReplicationNetworkID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationNetworkName string) ReplicationNetworkId {
	return ReplicationNetworkId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		VaultName:              vaultName,
		ReplicationFabricName:  replicationFabricName,
		ReplicationNetworkName: replicationNetworkName,
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
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.ReplicationFabricName, ok = parsed.Parsed["replicationFabricName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationFabricName", *parsed)
	}

	if id.ReplicationNetworkName, ok = parsed.Parsed["replicationNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationNetworkName", *parsed)
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
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.ReplicationFabricName, ok = parsed.Parsed["replicationFabricName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationFabricName", *parsed)
	}

	if id.ReplicationNetworkName, ok = parsed.Parsed["replicationNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationNetworkName", *parsed)
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationNetworkName)
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricValue"),
		resourceids.StaticSegment("staticReplicationNetworks", "replicationNetworks", "replicationNetworks"),
		resourceids.UserSpecifiedSegment("replicationNetworkName", "replicationNetworkValue"),
	}
}

// String returns a human-readable description of this Replication Network ID
func (id ReplicationNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Network Name: %q", id.ReplicationNetworkName),
	}
	return fmt.Sprintf("Replication Network (%s)", strings.Join(components, "\n"))
}
