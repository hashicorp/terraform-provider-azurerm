package iscsitargets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DiskPoolId{}

// DiskPoolId is a struct representing the Resource ID for a Disk Pool
type DiskPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	DiskPoolName      string
}

// NewDiskPoolID returns a new DiskPoolId struct
func NewDiskPoolID(subscriptionId string, resourceGroupName string, diskPoolName string) DiskPoolId {
	return DiskPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DiskPoolName:      diskPoolName,
	}
}

// ParseDiskPoolID parses 'input' into a DiskPoolId
func ParseDiskPoolID(input string) (*DiskPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(DiskPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DiskPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DiskPoolName, ok = parsed.Parsed["diskPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diskPoolName", *parsed)
	}

	return &id, nil
}

// ParseDiskPoolIDInsensitively parses 'input' case-insensitively into a DiskPoolId
// note: this method should only be used for API response data and not user input
func ParseDiskPoolIDInsensitively(input string) (*DiskPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(DiskPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DiskPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DiskPoolName, ok = parsed.Parsed["diskPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "diskPoolName", *parsed)
	}

	return &id, nil
}

// ValidateDiskPoolID checks that 'input' can be parsed as a Disk Pool ID
func ValidateDiskPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDiskPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disk Pool ID
func (id DiskPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StoragePool/diskPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DiskPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disk Pool ID
func (id DiskPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStoragePool", "Microsoft.StoragePool", "Microsoft.StoragePool"),
		resourceids.StaticSegment("staticDiskPools", "diskPools", "diskPools"),
		resourceids.UserSpecifiedSegment("diskPoolName", "diskPoolValue"),
	}
}

// String returns a human-readable description of this Disk Pool ID
func (id DiskPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Disk Pool Name: %q", id.DiskPoolName),
	}
	return fmt.Sprintf("Disk Pool (%s)", strings.Join(components, "\n"))
}
