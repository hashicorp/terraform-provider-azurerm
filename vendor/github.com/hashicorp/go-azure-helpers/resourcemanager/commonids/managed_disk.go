// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ManagedDiskId{}

// ManagedDiskId is a struct representing the Resource ID for a Managed Disk
type ManagedDiskId struct {
	SubscriptionId    string
	ResourceGroupName string
	DiskName          string
}

// NewManagedDiskID returns a new ManagedDiskId struct
func NewManagedDiskID(subscriptionId string, resourceGroupName string, diskName string) ManagedDiskId {
	return ManagedDiskId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DiskName:          diskName,
	}
}

// ParseManagedDiskID parses 'input' into a ManagedDiskId
func ParseManagedDiskID(input string) (*ManagedDiskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedDiskId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedDiskId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedDiskIDInsensitively parses 'input' case-insensitively into a ManagedDiskId
// note: this method should only be used for API response data and not user input
func ParseManagedDiskIDInsensitively(input string) (*ManagedDiskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedDiskId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedDiskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedDiskId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DiskName, ok = input.Parsed["diskName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "diskName", input)
	}

	return nil
}

// ValidateManagedDiskID checks that 'input' can be parsed as a Managed Disk ID
func ValidateManagedDiskID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedDiskID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Disk ID
func (id ManagedDiskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/disks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DiskName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disk ID
func (id ManagedDiskId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Managed Disk ID
func (id ManagedDiskId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Disk Name: %q", id.DiskName),
	}
	return fmt.Sprintf("Managed Disk (%s)", strings.Join(components, "\n"))
}
