package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SlotBackupId{})
}

var _ resourceids.ResourceId = &SlotBackupId{}

// SlotBackupId is a struct representing the Resource ID for a Slot Backup
type SlotBackupId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	BackupId          string
}

// NewSlotBackupID returns a new SlotBackupId struct
func NewSlotBackupID(subscriptionId string, resourceGroupName string, siteName string, slotName string, backupId string) SlotBackupId {
	return SlotBackupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		BackupId:          backupId,
	}
}

// ParseSlotBackupID parses 'input' into a SlotBackupId
func ParseSlotBackupID(input string) (*SlotBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotBackupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotBackupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotBackupIDInsensitively parses 'input' case-insensitively into a SlotBackupId
// note: this method should only be used for API response data and not user input
func ParseSlotBackupIDInsensitively(input string) (*SlotBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotBackupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotBackupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotBackupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.BackupId, ok = input.Parsed["backupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backupId", input)
	}

	return nil
}

// ValidateSlotBackupID checks that 'input' can be parsed as a Slot Backup ID
func ValidateSlotBackupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotBackupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Backup ID
func (id SlotBackupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/backups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.BackupId)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Backup ID
func (id SlotBackupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotName"),
		resourceids.StaticSegment("staticBackups", "backups", "backups"),
		resourceids.UserSpecifiedSegment("backupId", "backupId"),
	}
}

// String returns a human-readable description of this Slot Backup ID
func (id SlotBackupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Backup: %q", id.BackupId),
	}
	return fmt.Sprintf("Slot Backup (%s)", strings.Join(components, "\n"))
}
