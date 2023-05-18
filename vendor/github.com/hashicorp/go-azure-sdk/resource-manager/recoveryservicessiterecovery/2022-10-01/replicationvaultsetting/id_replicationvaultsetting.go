package replicationvaultsetting

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ReplicationVaultSettingId{}

// ReplicationVaultSettingId is a struct representing the Resource ID for a Replication Vault Setting
type ReplicationVaultSettingId struct {
	SubscriptionId              string
	ResourceGroupName           string
	VaultName                   string
	ReplicationVaultSettingName string
}

// NewReplicationVaultSettingID returns a new ReplicationVaultSettingId struct
func NewReplicationVaultSettingID(subscriptionId string, resourceGroupName string, vaultName string, replicationVaultSettingName string) ReplicationVaultSettingId {
	return ReplicationVaultSettingId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		VaultName:                   vaultName,
		ReplicationVaultSettingName: replicationVaultSettingName,
	}
}

// ParseReplicationVaultSettingID parses 'input' into a ReplicationVaultSettingId
func ParseReplicationVaultSettingID(input string) (*ReplicationVaultSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationVaultSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationVaultSettingId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.ReplicationVaultSettingName, ok = parsed.Parsed["replicationVaultSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationVaultSettingName", *parsed)
	}

	return &id, nil
}

// ParseReplicationVaultSettingIDInsensitively parses 'input' case-insensitively into a ReplicationVaultSettingId
// note: this method should only be used for API response data and not user input
func ParseReplicationVaultSettingIDInsensitively(input string) (*ReplicationVaultSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(ReplicationVaultSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ReplicationVaultSettingId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.ReplicationVaultSettingName, ok = parsed.Parsed["replicationVaultSettingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "replicationVaultSettingName", *parsed)
	}

	return &id, nil
}

// ValidateReplicationVaultSettingID checks that 'input' can be parsed as a Replication Vault Setting ID
func ValidateReplicationVaultSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseReplicationVaultSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Replication Vault Setting ID
func (id ReplicationVaultSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/replicationVaultSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationVaultSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Replication Vault Setting ID
func (id ReplicationVaultSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticReplicationVaultSettings", "replicationVaultSettings", "replicationVaultSettings"),
		resourceids.UserSpecifiedSegment("replicationVaultSettingName", "replicationVaultSettingValue"),
	}
}

// String returns a human-readable description of this Replication Vault Setting ID
func (id ReplicationVaultSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Replication Vault Setting Name: %q", id.ReplicationVaultSettingName),
	}
	return fmt.Sprintf("Replication Vault Setting (%s)", strings.Join(components, "\n"))
}
