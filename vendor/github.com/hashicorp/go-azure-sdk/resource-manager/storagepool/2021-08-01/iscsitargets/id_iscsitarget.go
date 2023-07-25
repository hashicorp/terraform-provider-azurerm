package iscsitargets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = IscsiTargetId{}

// IscsiTargetId is a struct representing the Resource ID for a Iscsi Target
type IscsiTargetId struct {
	SubscriptionId    string
	ResourceGroupName string
	DiskPoolName      string
	IscsiTargetName   string
}

// NewIscsiTargetID returns a new IscsiTargetId struct
func NewIscsiTargetID(subscriptionId string, resourceGroupName string, diskPoolName string, iscsiTargetName string) IscsiTargetId {
	return IscsiTargetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DiskPoolName:      diskPoolName,
		IscsiTargetName:   iscsiTargetName,
	}
}

// ParseIscsiTargetID parses 'input' into a IscsiTargetId
func ParseIscsiTargetID(input string) (*IscsiTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(IscsiTargetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IscsiTargetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DiskPoolName, ok = parsed.Parsed["diskPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diskPoolName", *parsed)
	}

	if id.IscsiTargetName, ok = parsed.Parsed["iscsiTargetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "iscsiTargetName", *parsed)
	}

	return &id, nil
}

// ParseIscsiTargetIDInsensitively parses 'input' case-insensitively into a IscsiTargetId
// note: this method should only be used for API response data and not user input
func ParseIscsiTargetIDInsensitively(input string) (*IscsiTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(IscsiTargetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IscsiTargetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DiskPoolName, ok = parsed.Parsed["diskPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diskPoolName", *parsed)
	}

	if id.IscsiTargetName, ok = parsed.Parsed["iscsiTargetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "iscsiTargetName", *parsed)
	}

	return &id, nil
}

// ValidateIscsiTargetID checks that 'input' can be parsed as a Iscsi Target ID
func ValidateIscsiTargetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIscsiTargetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Iscsi Target ID
func (id IscsiTargetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StoragePool/diskPools/%s/iscsiTargets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DiskPoolName, id.IscsiTargetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Iscsi Target ID
func (id IscsiTargetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStoragePool", "Microsoft.StoragePool", "Microsoft.StoragePool"),
		resourceids.StaticSegment("staticDiskPools", "diskPools", "diskPools"),
		resourceids.UserSpecifiedSegment("diskPoolName", "diskPoolValue"),
		resourceids.StaticSegment("staticIscsiTargets", "iscsiTargets", "iscsiTargets"),
		resourceids.UserSpecifiedSegment("iscsiTargetName", "iscsiTargetValue"),
	}
}

// String returns a human-readable description of this Iscsi Target ID
func (id IscsiTargetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Disk Pool Name: %q", id.DiskPoolName),
		fmt.Sprintf("Iscsi Target Name: %q", id.IscsiTargetName),
	}
	return fmt.Sprintf("Iscsi Target (%s)", strings.Join(components, "\n"))
}
