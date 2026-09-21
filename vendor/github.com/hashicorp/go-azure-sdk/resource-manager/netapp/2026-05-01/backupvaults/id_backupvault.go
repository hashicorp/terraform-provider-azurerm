package backupvaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackupVaultId{})
}

var _ resourceids.ResourceId = &BackupVaultId{}

// BackupVaultId is a struct representing the Resource ID for a Backup Vault
type BackupVaultId struct {
	SubscriptionId    string
	ResourceGroupName string
	NetAppAccountName string
	BackupVaultName   string
}

// NewBackupVaultID returns a new BackupVaultId struct
func NewBackupVaultID(subscriptionId string, resourceGroupName string, netAppAccountName string, backupVaultName string) BackupVaultId {
	return BackupVaultId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NetAppAccountName: netAppAccountName,
		BackupVaultName:   backupVaultName,
	}
}

// ParseBackupVaultID parses 'input' into a BackupVaultId
func ParseBackupVaultID(input string) (*BackupVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupVaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupVaultIDInsensitively parses 'input' case-insensitively into a BackupVaultId
// note: this method should only be used for API response data and not user input
func ParseBackupVaultIDInsensitively(input string) (*BackupVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupVaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupVaultId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetAppAccountName, ok = input.Parsed["netAppAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "netAppAccountName", input)
	}

	if id.BackupVaultName, ok = input.Parsed["backupVaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupVaultName", input)
	}

	return nil
}

// ValidateBackupVaultID checks that 'input' can be parsed as a Backup Vault ID
func ValidateBackupVaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Vault ID
func (id BackupVaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/backupVaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.BackupVaultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Vault ID
func (id BackupVaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("netAppAccountName", "netAppAccountName"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("backupVaultName", "backupVaultName"),
	}
}

// String returns a human-readable description of this Backup Vault ID
func (id BackupVaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Net App Account Name: %q", id.NetAppAccountName),
		fmt.Sprintf("Backup Vault Name: %q", id.BackupVaultName),
	}
	return fmt.Sprintf("Backup Vault (%s)", strings.Join(components, "\n"))
}
