package syncgroupresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = StorageSyncServiceId{}

// StorageSyncServiceId is a struct representing the Resource ID for a Storage Sync Service
type StorageSyncServiceId struct {
	SubscriptionId         string
	ResourceGroupName      string
	StorageSyncServiceName string
}

// NewStorageSyncServiceID returns a new StorageSyncServiceId struct
func NewStorageSyncServiceID(subscriptionId string, resourceGroupName string, storageSyncServiceName string) StorageSyncServiceId {
	return StorageSyncServiceId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		StorageSyncServiceName: storageSyncServiceName,
	}
}

// ParseStorageSyncServiceID parses 'input' into a StorageSyncServiceId
func ParseStorageSyncServiceID(input string) (*StorageSyncServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageSyncServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageSyncServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageSyncServiceName, ok = parsed.Parsed["storageSyncServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageSyncServiceName", *parsed)
	}

	return &id, nil
}

// ParseStorageSyncServiceIDInsensitively parses 'input' case-insensitively into a StorageSyncServiceId
// note: this method should only be used for API response data and not user input
func ParseStorageSyncServiceIDInsensitively(input string) (*StorageSyncServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageSyncServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageSyncServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageSyncServiceName, ok = parsed.Parsed["storageSyncServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageSyncServiceName", *parsed)
	}

	return &id, nil
}

// ValidateStorageSyncServiceID checks that 'input' can be parsed as a Storage Sync Service ID
func ValidateStorageSyncServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageSyncServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Sync Service ID
func (id StorageSyncServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageSync/storageSyncServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Sync Service ID
func (id StorageSyncServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageSync", "Microsoft.StorageSync", "Microsoft.StorageSync"),
		resourceids.StaticSegment("staticStorageSyncServices", "storageSyncServices", "storageSyncServices"),
		resourceids.UserSpecifiedSegment("storageSyncServiceName", "storageSyncServiceValue"),
	}
}

// String returns a human-readable description of this Storage Sync Service ID
func (id StorageSyncServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Sync Service Name: %q", id.StorageSyncServiceName),
	}
	return fmt.Sprintf("Storage Sync Service (%s)", strings.Join(components, "\n"))
}
