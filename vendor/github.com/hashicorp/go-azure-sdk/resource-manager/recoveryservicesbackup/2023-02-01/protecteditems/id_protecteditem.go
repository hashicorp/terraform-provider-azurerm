package protecteditems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProtectedItemId{}

// ProtectedItemId is a struct representing the Resource ID for a Protected Item
type ProtectedItemId struct {
	SubscriptionId          string
	ResourceGroupName       string
	VaultName               string
	BackupFabricName        string
	ProtectionContainerName string
	ProtectedItemName       string
}

// NewProtectedItemID returns a new ProtectedItemId struct
func NewProtectedItemID(subscriptionId string, resourceGroupName string, vaultName string, backupFabricName string, protectionContainerName string, protectedItemName string) ProtectedItemId {
	return ProtectedItemId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		VaultName:               vaultName,
		BackupFabricName:        backupFabricName,
		ProtectionContainerName: protectionContainerName,
		ProtectedItemName:       protectedItemName,
	}
}

// ParseProtectedItemID parses 'input' into a ProtectedItemId
func ParseProtectedItemID(input string) (*ProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProtectedItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProtectedItemId{}

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

	if id.ProtectedItemName, ok = parsed.Parsed["protectedItemName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "protectedItemName", *parsed)
	}

	return &id, nil
}

// ParseProtectedItemIDInsensitively parses 'input' case-insensitively into a ProtectedItemId
// note: this method should only be used for API response data and not user input
func ParseProtectedItemIDInsensitively(input string) (*ProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProtectedItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProtectedItemId{}

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

	if id.ProtectedItemName, ok = parsed.Parsed["protectedItemName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "protectedItemName", *parsed)
	}

	return &id, nil
}

// ValidateProtectedItemID checks that 'input' can be parsed as a Protected Item ID
func ValidateProtectedItemID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProtectedItemID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Protected Item ID
func (id ProtectedItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/backupFabrics/%s/protectionContainers/%s/protectedItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.BackupFabricName, id.ProtectionContainerName, id.ProtectedItemName)
}

// Segments returns a slice of Resource ID Segments which comprise this Protected Item ID
func (id ProtectedItemId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticProtectedItems", "protectedItems", "protectedItems"),
		resourceids.UserSpecifiedSegment("protectedItemName", "protectedItemValue"),
	}
}

// String returns a human-readable description of this Protected Item ID
func (id ProtectedItemId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Backup Fabric Name: %q", id.BackupFabricName),
		fmt.Sprintf("Protection Container Name: %q", id.ProtectionContainerName),
		fmt.Sprintf("Protected Item Name: %q", id.ProtectedItemName),
	}
	return fmt.Sprintf("Protected Item (%s)", strings.Join(components, "\n"))
}
