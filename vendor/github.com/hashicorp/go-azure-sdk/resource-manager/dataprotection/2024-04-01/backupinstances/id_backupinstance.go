package backupinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackupInstanceId{})
}

var _ resourceids.ResourceId = &BackupInstanceId{}

// BackupInstanceId is a struct representing the Resource ID for a Backup Instance
type BackupInstanceId struct {
	SubscriptionId     string
	ResourceGroupName  string
	BackupVaultName    string
	BackupInstanceName string
}

// NewBackupInstanceID returns a new BackupInstanceId struct
func NewBackupInstanceID(subscriptionId string, resourceGroupName string, backupVaultName string, backupInstanceName string) BackupInstanceId {
	return BackupInstanceId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		BackupVaultName:    backupVaultName,
		BackupInstanceName: backupInstanceName,
	}
}

// ParseBackupInstanceID parses 'input' into a BackupInstanceId
func ParseBackupInstanceID(input string) (*BackupInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupInstanceIDInsensitively parses 'input' case-insensitively into a BackupInstanceId
// note: this method should only be used for API response data and not user input
func ParseBackupInstanceIDInsensitively(input string) (*BackupInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.BackupVaultName, ok = input.Parsed["backupVaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupVaultName", input)
	}

	if id.BackupInstanceName, ok = input.Parsed["backupInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupInstanceName", input)
	}

	return nil
}

// ValidateBackupInstanceID checks that 'input' can be parsed as a Backup Instance ID
func ValidateBackupInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Instance ID
func (id BackupInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s/backupInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName, id.BackupInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Instance ID
func (id BackupInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("backupVaultName", "backupVaultName"),
		resourceids.StaticSegment("staticBackupInstances", "backupInstances", "backupInstances"),
		resourceids.UserSpecifiedSegment("backupInstanceName", "backupInstanceName"),
	}
}

// String returns a human-readable description of this Backup Instance ID
func (id BackupInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Backup Vault Name: %q", id.BackupVaultName),
		fmt.Sprintf("Backup Instance Name: %q", id.BackupInstanceName),
	}
	return fmt.Sprintf("Backup Instance (%s)", strings.Join(components, "\n"))
}
