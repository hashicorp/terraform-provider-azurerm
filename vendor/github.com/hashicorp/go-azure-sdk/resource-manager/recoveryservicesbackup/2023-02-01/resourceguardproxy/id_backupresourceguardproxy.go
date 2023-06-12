package resourceguardproxy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BackupResourceGuardProxyId{}

// BackupResourceGuardProxyId is a struct representing the Resource ID for a Backup Resource Guard Proxy
type BackupResourceGuardProxyId struct {
	SubscriptionId               string
	ResourceGroupName            string
	VaultName                    string
	BackupResourceGuardProxyName string
}

// NewBackupResourceGuardProxyID returns a new BackupResourceGuardProxyId struct
func NewBackupResourceGuardProxyID(subscriptionId string, resourceGroupName string, vaultName string, backupResourceGuardProxyName string) BackupResourceGuardProxyId {
	return BackupResourceGuardProxyId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		VaultName:                    vaultName,
		BackupResourceGuardProxyName: backupResourceGuardProxyName,
	}
}

// ParseBackupResourceGuardProxyID parses 'input' into a BackupResourceGuardProxyId
func ParseBackupResourceGuardProxyID(input string) (*BackupResourceGuardProxyId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupResourceGuardProxyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupResourceGuardProxyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.BackupResourceGuardProxyName, ok = parsed.Parsed["backupResourceGuardProxyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "backupResourceGuardProxyName", *parsed)
	}

	return &id, nil
}

// ParseBackupResourceGuardProxyIDInsensitively parses 'input' case-insensitively into a BackupResourceGuardProxyId
// note: this method should only be used for API response data and not user input
func ParseBackupResourceGuardProxyIDInsensitively(input string) (*BackupResourceGuardProxyId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupResourceGuardProxyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupResourceGuardProxyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.BackupResourceGuardProxyName, ok = parsed.Parsed["backupResourceGuardProxyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "backupResourceGuardProxyName", *parsed)
	}

	return &id, nil
}

// ValidateBackupResourceGuardProxyID checks that 'input' can be parsed as a Backup Resource Guard Proxy ID
func ValidateBackupResourceGuardProxyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupResourceGuardProxyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Resource Guard Proxy ID
func (id BackupResourceGuardProxyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/backupResourceGuardProxies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.BackupResourceGuardProxyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Resource Guard Proxy ID
func (id BackupResourceGuardProxyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticBackupResourceGuardProxies", "backupResourceGuardProxies", "backupResourceGuardProxies"),
		resourceids.UserSpecifiedSegment("backupResourceGuardProxyName", "backupResourceGuardProxyValue"),
	}
}

// String returns a human-readable description of this Backup Resource Guard Proxy ID
func (id BackupResourceGuardProxyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Backup Resource Guard Proxy Name: %q", id.BackupResourceGuardProxyName),
	}
	return fmt.Sprintf("Backup Resource Guard Proxy (%s)", strings.Join(components, "\n"))
}
