package replicationprotectioncontainermappings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicationProtectionContainerMappingId{})
}

var _ resourceids.ResourceId = &ReplicationProtectionContainerMappingId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ReplicationProtectionContainerMappingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationProtectionContainerMappingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationProtectionContainerMappingIDInsensitively parses 'input' case-insensitively into a ReplicationProtectionContainerMappingId
// note: this method should only be used for API response data and not user input
func ParseReplicationProtectionContainerMappingIDInsensitively(input string) (*ReplicationProtectionContainerMappingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationProtectionContainerMappingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationProtectionContainerMappingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationProtectionContainerMappingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VaultName, ok = input.Parsed["vaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vaultName", input)
	}

	if id.ReplicationFabricName, ok = input.Parsed["replicationFabricName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationFabricName", input)
	}

	if id.ReplicationProtectionContainerName, ok = input.Parsed["replicationProtectionContainerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectionContainerName", input)
	}

	if id.ReplicationProtectionContainerMappingName, ok = input.Parsed["replicationProtectionContainerMappingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectionContainerMappingName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricName"),
		resourceids.StaticSegment("staticReplicationProtectionContainers", "replicationProtectionContainers", "replicationProtectionContainers"),
		resourceids.UserSpecifiedSegment("replicationProtectionContainerName", "replicationProtectionContainerName"),
		resourceids.StaticSegment("staticReplicationProtectionContainerMappings", "replicationProtectionContainerMappings", "replicationProtectionContainerMappings"),
		resourceids.UserSpecifiedSegment("replicationProtectionContainerMappingName", "replicationProtectionContainerMappingName"),
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
