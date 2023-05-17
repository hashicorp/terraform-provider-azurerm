package replicationprotectioncontainermappings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationProtectionContainerMappingId{}

// ReplicationProtectionContainerMappingId is a struct representing the Resource ID for a Replication Protection Container Mapping
type ReplicationProtectionContainerMappingId struct {
	SubscriptionId                            string
	ResourceGroupName                         string
	VaultName                                 string
	ReplicationFabricName                     string
	ReplicationProtectionContainerName        string
	ReplicationProtectionContainerMappingName string
}

// NewReplicationProtectionContainerMappingID returns a new ReplicationProtectionContainerMappingId struct
func NewReplicationProtectionContainerMappingID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationProtectionContainerName string, replicationProtectionContainerMappingName string) ReplicationProtectionContainerMappingId {
	return ReplicationProtectionContainerMappingId{
		SubscriptionId:                            subscriptionId,
		ResourceGroupName:                         resourceGroupName,
		VaultName:                                 vaultName,
		ReplicationFabricName:                     replicationFabricName,
		ReplicationProtectionContainerName:        replicationProtectionContainerName,
		ReplicationProtectionContainerMappingName: replicationProtectionContainerMappingName,
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

	if id.ReplicationProtectionContainerMappingName, ok = parsed.Parsed["replicationProtectionContainerMappingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectionContainerMappingName", *parsed)
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

	if id.ReplicationProtectionContainerMappingName, ok = parsed.Parsed["replicationProtectionContainerMappingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectionContainerMappingName", *parsed)
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectionContainerMappingName)
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricValue"),
		resourceids.StaticSegment("staticReplicationProtectionContainers", "replicationProtectionContainers", "replicationProtectionContainers"),
		resourceids.UserSpecifiedSegment("replicationProtectionContainerName", "replicationProtectionContainerValue"),
		resourceids.StaticSegment("staticReplicationProtectionContainerMappings", "replicationProtectionContainerMappings", "replicationProtectionContainerMappings"),
		resourceids.UserSpecifiedSegment("replicationProtectionContainerMappingName", "replicationProtectionContainerMappingValue"),
	}
}

// String returns a human-readable description of this Replication Protection Container Mapping ID
func (id ReplicationProtectionContainerMappingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Protection Container Name: %q", id.ReplicationProtectionContainerName),
		fmt.Sprintf("Replication Protection Container Mapping Name: %q", id.ReplicationProtectionContainerMappingName),
	}
	return fmt.Sprintf("Replication Protection Container Mapping (%s)", strings.Join(components, "\n"))
}
