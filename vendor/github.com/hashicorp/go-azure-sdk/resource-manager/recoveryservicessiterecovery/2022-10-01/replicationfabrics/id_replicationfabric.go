package replicationfabrics

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationFabricId{}

// ReplicationFabricId is a struct representing the Resource ID for a Replication Fabric
type ReplicationFabricId struct {
	SubscriptionId        string
	ResourceGroupName     string
	VaultName             string
	ReplicationFabricName string
}

// NewReplicationFabricID returns a new ReplicationFabricId struct
func NewReplicationFabricID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string) ReplicationFabricId {
	return ReplicationFabricId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		VaultName:             vaultName,
		ReplicationFabricName: replicationFabricName,
	}
}

// ParseReplicationFabricID parses 'input' into a ReplicationFabricId
func ParseReplicationFabricID(input string) (*ReplicationFabricId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationFabricId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationFabricId{}

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

	return &id, nil
}

// ParseReplicationFabricIDInsensitively parses 'input' case-insensitively into a ReplicationFabricId
// note: this method should only be used for API response data and not user input
func ParseReplicationFabricIDInsensitively(input string) (*ReplicationFabricId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationFabricId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationFabricId{}

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

	return &id, nil
}

// ValidateReplicationFabricID checks that 'input' can be parsed as a Replication Fabric ID
func ValidateReplicationFabricID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationFabricID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Fabric ID
func (id ReplicationFabricId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Fabric ID
func (id ReplicationFabricId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Replication Fabric ID
func (id ReplicationFabricId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
	}
	return fmt.Sprintf("Replication Fabric (%s)", strings.Join(components, "\n"))
}
