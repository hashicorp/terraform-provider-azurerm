package replicationprotecteditems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationProtectionContainerId{}

// ReplicationProtectionContainerId is a struct representing the Resource ID for a Replication Protection Container
type ReplicationProtectionContainerId struct {
	SubscriptionId                     string
	ResourceGroupName                  string
	VaultName                          string
	ReplicationFabricName              string
	ReplicationProtectionContainerName string
}

// NewReplicationProtectionContainerID returns a new ReplicationProtectionContainerId struct
func NewReplicationProtectionContainerID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationProtectionContainerName string) ReplicationProtectionContainerId {
	return ReplicationProtectionContainerId{
		SubscriptionId:                     subscriptionId,
		ResourceGroupName:                  resourceGroupName,
		VaultName:                          vaultName,
		ReplicationFabricName:              replicationFabricName,
		ReplicationProtectionContainerName: replicationProtectionContainerName,
	}
}

// ParseReplicationProtectionContainerID parses 'input' into a ReplicationProtectionContainerId
func ParseReplicationProtectionContainerID(input string) (*ReplicationProtectionContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectionContainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectionContainerId{}

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

	if id.ReplicationProtectionContainerName, ok = parsed.Parsed["replicationProtectionContainerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectionContainerName", *parsed)
	}

	return &id, nil
}

// ParseReplicationProtectionContainerIDInsensitively parses 'input' case-insensitively into a ReplicationProtectionContainerId
// note: this method should only be used for API response data and not user input
func ParseReplicationProtectionContainerIDInsensitively(input string) (*ReplicationProtectionContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationProtectionContainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationProtectionContainerId{}

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

	if id.ReplicationProtectionContainerName, ok = parsed.Parsed["replicationProtectionContainerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectionContainerName", *parsed)
	}

	return &id, nil
}

// ValidateReplicationProtectionContainerID checks that 'input' can be parsed as a Replication Protection Container ID
func ValidateReplicationProtectionContainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationProtectionContainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Protection Container ID
func (id ReplicationProtectionContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationProtectionContainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationProtectionContainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Protection Container ID
func (id ReplicationProtectionContainerId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Replication Protection Container ID
func (id ReplicationProtectionContainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Protection Container Name: %q", id.ReplicationProtectionContainerName),
	}
	return fmt.Sprintf("Replication Protection Container (%s)", strings.Join(components, "\n"))
}
