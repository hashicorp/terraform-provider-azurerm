package disks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DiskId{}

// DiskId is a struct representing the Resource ID for a Disk
type DiskId struct {
	SubscriptionId    string
	ResourceGroupName string
	DiskName          string
}

// NewDiskID returns a new DiskId struct
func NewDiskID(subscriptionId string, resourceGroupName string, diskName string) DiskId {
	return DiskId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DiskName:          diskName,
	}
}

// ParseDiskID parses 'input' into a DiskId
func ParseDiskID(input string) (*DiskId, error) {
	parser := resourceids.NewParserFromResourceIdType(DiskId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DiskId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DiskName, ok = parsed.Parsed["diskName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diskName", *parsed)
	}

	return &id, nil
}

// ParseDiskIDInsensitively parses 'input' case-insensitively into a DiskId
// note: this method should only be used for API response data and not user input
func ParseDiskIDInsensitively(input string) (*DiskId, error) {
	parser := resourceids.NewParserFromResourceIdType(DiskId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DiskId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DiskName, ok = parsed.Parsed["diskName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diskName", *parsed)
	}

	return &id, nil
}

// ValidateDiskID checks that 'input' can be parsed as a Disk ID
func ValidateDiskID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDiskID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disk ID
func (id DiskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/disks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DiskName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disk ID
func (id DiskId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticDisks", "disks", "disks"),
		resourceids.UserSpecifiedSegment("diskName", "diskValue"),
	}
}

// String returns a human-readable description of this Disk ID
func (id DiskId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Disk Name: %q", id.DiskName),
	}
	return fmt.Sprintf("Disk (%s)", strings.Join(components, "\n"))
}
