package protectioncontainers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackupFabricId{})
}

var _ resourceids.ResourceId = &BackupFabricId{}

// BackupFabricId is a struct representing the Resource ID for a Backup Fabric
type BackupFabricId struct {
	SubscriptionId    string
	ResourceGroupName string
	VaultName         string
	BackupFabricName  string
}

// NewBackupFabricID returns a new BackupFabricId struct
func NewBackupFabricID(subscriptionId string, resourceGroupName string, vaultName string, backupFabricName string) BackupFabricId {
	return BackupFabricId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VaultName:         vaultName,
		BackupFabricName:  backupFabricName,
	}
}

// ParseBackupFabricID parses 'input' into a BackupFabricId
func ParseBackupFabricID(input string) (*BackupFabricId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupFabricId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupFabricId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupFabricIDInsensitively parses 'input' case-insensitively into a BackupFabricId
// note: this method should only be used for API response data and not user input
func ParseBackupFabricIDInsensitively(input string) (*BackupFabricId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupFabricId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupFabricId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupFabricId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateBackupFabricID checks that 'input' can be parsed as a Backup Fabric ID
func ValidateBackupFabricID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupFabricID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Fabric ID
func (id BackupFabricId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/backupFabrics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.BackupFabricName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Fabric ID
func (id BackupFabricId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Backup Fabric ID
func (id BackupFabricId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Backup Fabric Name: %q", id.BackupFabricName),
	}
	return fmt.Sprintf("Backup Fabric (%s)", strings.Join(components, "\n"))
}
