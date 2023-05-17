package replicationrecoveryservicesproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationRecoveryServicesProviderId{}

// ReplicationRecoveryServicesProviderId is a struct representing the Resource ID for a Replication Recovery Services Provider
type ReplicationRecoveryServicesProviderId struct {
	SubscriptionId                          string
	ResourceGroupName                       string
	VaultName                               string
	ReplicationFabricName                   string
	ReplicationRecoveryServicesProviderName string
}

// NewReplicationRecoveryServicesProviderID returns a new ReplicationRecoveryServicesProviderId struct
func NewReplicationRecoveryServicesProviderID(subscriptionId string, resourceGroupName string, vaultName string, replicationFabricName string, replicationRecoveryServicesProviderName string) ReplicationRecoveryServicesProviderId {
	return ReplicationRecoveryServicesProviderId{
		SubscriptionId:                          subscriptionId,
		ResourceGroupName:                       resourceGroupName,
		VaultName:                               vaultName,
		ReplicationFabricName:                   replicationFabricName,
		ReplicationRecoveryServicesProviderName: replicationRecoveryServicesProviderName,
	}
}

// ParseReplicationRecoveryServicesProviderID parses 'input' into a ReplicationRecoveryServicesProviderId
func ParseReplicationRecoveryServicesProviderID(input string) (*ReplicationRecoveryServicesProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationRecoveryServicesProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationRecoveryServicesProviderId{}

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

	if id.ReplicationRecoveryServicesProviderName, ok = parsed.Parsed["replicationRecoveryServicesProviderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationRecoveryServicesProviderName", *parsed)
	}

	return &id, nil
}

// ParseReplicationRecoveryServicesProviderIDInsensitively parses 'input' case-insensitively into a ReplicationRecoveryServicesProviderId
// note: this method should only be used for API response data and not user input
func ParseReplicationRecoveryServicesProviderIDInsensitively(input string) (*ReplicationRecoveryServicesProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationRecoveryServicesProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationRecoveryServicesProviderId{}

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

	if id.ReplicationRecoveryServicesProviderName, ok = parsed.Parsed["replicationRecoveryServicesProviderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationRecoveryServicesProviderName", *parsed)
	}

	return &id, nil
}

// ValidateReplicationRecoveryServicesProviderID checks that 'input' can be parsed as a Replication Recovery Services Provider ID
func ValidateReplicationRecoveryServicesProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationRecoveryServicesProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Recovery Services Provider ID
func (id ReplicationRecoveryServicesProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationFabrics/%s/replicationRecoveryServicesProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName, id.ReplicationRecoveryServicesProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Recovery Services Provider ID
func (id ReplicationRecoveryServicesProviderId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticReplicationRecoveryServicesProviders", "replicationRecoveryServicesProviders", "replicationRecoveryServicesProviders"),
		resourceids.UserSpecifiedSegment("replicationRecoveryServicesProviderName", "replicationRecoveryServicesProviderValue"),
	}
}

// String returns a human-readable description of this Replication Recovery Services Provider ID
func (id ReplicationRecoveryServicesProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Fabric Name: %q", id.ReplicationFabricName),
		fmt.Sprintf("Replication Recovery Services Provider Name: %q", id.ReplicationRecoveryServicesProviderName),
	}
	return fmt.Sprintf("Replication Recovery Services Provider (%s)", strings.Join(components, "\n"))
}
