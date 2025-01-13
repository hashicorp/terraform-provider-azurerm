package replicationvaultsetting

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ReplicationVaultSettingId{})
}

var _ resourceids.ResourceId = &ReplicationVaultSettingId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ReplicationVaultSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationVaultSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseReplicationVaultSettingIDInsensitively parses 'input' case-insensitively into a ReplicationVaultSettingId
// note: this method should only be used for API response data and not user input
func ParseReplicationVaultSettingIDInsensitively(input string) (*ReplicationVaultSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ReplicationVaultSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ReplicationVaultSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ReplicationVaultSettingId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ReplicationVaultSettingName, ok = input.Parsed["replicationVaultSettingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "replicationVaultSettingName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticReplicationVaultSettings", "replicationVaultSettings", "replicationVaultSettings"),
		resourceids.UserSpecifiedSegment("replicationVaultSettingName", "replicationVaultSettingName"),
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
