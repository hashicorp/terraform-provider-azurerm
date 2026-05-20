package replicationprotecteditems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicationProtectedItemId{})
}

var _ resourceids.ResourceId = &ReplicationProtectedItemId{}

// ReplicationProtectedItemId is a struct representing the Resource ID for a Replication Protected Item
type ReplicationProtectedItemId struct {
	SubscriptionId                     string
	ResourceGroupName                  string
	VaultName                          string
	ReplicationFabricName              string
	ReplicationProtectionContainerName string
	ReplicationProtectedItemName       string
}

// NewReplicationProtectedItemID returns a new ReplicationProtectedItemId struct
func NewReplicationProtectedItemID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationProtectionContainerName string, replicationProtectedItemName string) ReplicationProtectedItemId {
	return ReplicationProtectedItemId{
		SubscriptionId:                     subscriptionId,
		ResourceGroupName:                  resourceGroupName,
		VaultName:                          vaultName,
		ReplicationFabricName:              replicationFabricName,
		ReplicationProtectionContainerName: replicationProtectionContainerName,
		ReplicationProtectedItemName:       replicationProtectedItemName,
	}
}

// ParseReplicationProtectedItemID parses 'input' into a ReplicationProtectedItemId
func ParseReplicationProtectedItemID(input string) (*ReplicationProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationProtectedItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationProtectedItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationProtectedItemIDInsensitively parses 'input' case-insensitively into a ReplicationProtectedItemId
// note: this method should only be used for API response data and not user input
func ParseReplicationProtectedItemIDInsensitively(input string) (*ReplicationProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationProtectedItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationProtectedItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationProtectedItemId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ReplicationProtectedItemName, ok = input.Parsed["replicationProtectedItemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationProtectedItemName", input)
	}

	return nil
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectedItemName)
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricName"),
		resourceids.StaticSegment("staticReplicationProtectionContainers", "replicationProtectionContainers", "replicationProtectionContainers"),
		resourceids.UserSpecifiedSegment("replicationProtectionContainerName", "replicationProtectionContainerName"),
		resourceids.StaticSegment("staticReplicationProtectedItems", "replicationProtectedItems", "replicationProtectedItems"),
		resourceids.UserSpecifiedSegment("replicationProtectedItemName", "replicationProtectedItemName"),
	}
}

// String returns a human-readable description of this Replication Protected Item ID
func (id ReplicationProtectedItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Protection Container Name: %q", id.ReplicationProtectionContainerName),
		fmt.Sprintf("Replication Protected Item Name: %q", id.ReplicationProtectedItemName),
	}
	return fmt.Sprintf("Replication Protected Item (%s)", strings.Join(components, "\n"))
}
