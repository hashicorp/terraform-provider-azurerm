package protecteditems

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProtectedItemId{})
}

var _ resourceids.ResourceId = &ProtectedItemId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ProtectedItemId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProtectedItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProtectedItemIDInsensitively parses 'input' case-insensitively into a ProtectedItemId
// note: this method should only be used for API response data and not user input
func ParseProtectedItemIDInsensitively(input string) (*ProtectedItemId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProtectedItemId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProtectedItemId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProtectedItemId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BackupFabricName, ok = input.Parsed["backupFabricName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupFabricName", input)
	}

	if id.ProtectionContainerName, ok = input.Parsed["protectionContainerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "protectionContainerName", input)
	}

	if id.ProtectedItemName, ok = input.Parsed["protectedItemName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "protectedItemName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticBackupFabrics", "backupFabrics", "backupFabrics"),
		resourceids.UserSpecifiedSegment("backupFabricName", "backupFabricName"),
		resourceids.StaticSegment("staticProtectionContainers", "protectionContainers", "protectionContainers"),
		resourceids.UserSpecifiedSegment("protectionContainerName", "protectionContainerName"),
		resourceids.StaticSegment("staticProtectedItems", "protectedItems", "protectedItems"),
		resourceids.UserSpecifiedSegment("protectedItemName", "protectedItemName"),
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
