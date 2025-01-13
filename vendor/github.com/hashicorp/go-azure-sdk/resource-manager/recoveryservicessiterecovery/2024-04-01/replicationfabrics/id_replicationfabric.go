package replicationfabrics

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicationFabricId{})
}

var _ resourceids.ResourceId = &ReplicationFabricId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ReplicationFabricId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationFabricId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationFabricIDInsensitively parses 'input' case-insensitively into a ReplicationFabricId
// note: this method should only be used for API response data and not user input
func ParseReplicationFabricIDInsensitively(input string) (*ReplicationFabricId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationFabricId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationFabricId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationFabricId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticReplicationFabrics", "replicationFabrics", "replicationFabrics"),
		resourceids.UserSpecifiedSegment("replicationFabricName", "replicationFabricName"),
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
