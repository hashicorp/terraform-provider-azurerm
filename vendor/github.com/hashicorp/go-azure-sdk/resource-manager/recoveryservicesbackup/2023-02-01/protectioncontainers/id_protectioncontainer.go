package protectioncontainers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProtectionContainerId{}

// ProtectionContainerId is a struct representing the Resource ID for a Protection Container
type ProtectionContainerId struct {
	SubscriptionId          string
	ResourceGroupName       string
	VaultName               string
	BackupFabricName        string
	ProtectionContainerName string
}

// NewProtectionContainerID returns a new ProtectionContainerId struct
func NewProtectionContainerID(subscriptionId string, resourceGroupName string, vaultName string, backupFabricName string, protectionContainerName string) ProtectionContainerId {
	return ProtectionContainerId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		VaultName:               vaultName,
		BackupFabricName:        backupFabricName,
		ProtectionContainerName: protectionContainerName,
	}
}

// ParseProtectionContainerID parses 'input' into a ProtectionContainerId
func ParseProtectionContainerID(input string) (*ProtectionContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProtectionContainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProtectionContainerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.BackupFabricName, ok = parsed.Parsed["backupFabricName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "backupFabricName", *parsed)
	}

	if id.ProtectionContainerName, ok = parsed.Parsed["protectionContainerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "protectionContainerName", *parsed)
	}

	return &id, nil
}

// ParseProtectionContainerIDInsensitively parses 'input' case-insensitively into a ProtectionContainerId
// note: this method should only be used for API response data and not user input
func ParseProtectionContainerIDInsensitively(input string) (*ProtectionContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProtectionContainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProtectionContainerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vaultName", *parsed)
	}

	if id.BackupFabricName, ok = parsed.Parsed["backupFabricName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "backupFabricName", *parsed)
	}

	if id.ProtectionContainerName, ok = parsed.Parsed["protectionContainerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "protectionContainerName", *parsed)
	}

	return &id, nil
}

// ValidateProtectionContainerID checks that 'input' can be parsed as a Protection Container ID
func ValidateProtectionContainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProtectionContainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Protection Container ID
func (id ProtectionContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/backupFabrics/%s/protectionContainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.BackupFabricName, id.ProtectionContainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Protection Container ID
func (id ProtectionContainerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticBackupFabrics", "backupFabrics", "backupFabrics"),
		resourceids.UserSpecifiedSegment("backupFabricName", "backupFabricValue"),
		resourceids.StaticSegment("staticProtectionContainers", "protectionContainers", "protectionContainers"),
		resourceids.UserSpecifiedSegment("protectionContainerName", "protectionContainerValue"),
	}
}

// String returns a human-readable description of this Protection Container ID
func (id ProtectionContainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Backup Fabric Name: %q", id.BackupFabricName),
		fmt.Sprintf("Protection Container Name: %q", id.ProtectionContainerName),
	}
	return fmt.Sprintf("Protection Container (%s)", strings.Join(components, "\n"))
}
