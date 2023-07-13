package projects

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = StorageMoverId{}

// StorageMoverId is a struct representing the Resource ID for a Storage Mover
type StorageMoverId struct {
	SubscriptionId    string
	ResourceGroupName string
	StorageMoverName  string
}

// NewStorageMoverID returns a new StorageMoverId struct
func NewStorageMoverID(subscriptionId string, resourceGroupName string, storageMoverName string) StorageMoverId {
	return StorageMoverId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StorageMoverName:  storageMoverName,
	}
}

// ParseStorageMoverID parses 'input' into a StorageMoverId
func ParseStorageMoverID(input string) (*StorageMoverId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageMoverId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageMoverId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageMoverName, ok = parsed.Parsed["storageMoverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", *parsed)
	}

	return &id, nil
}

// ParseStorageMoverIDInsensitively parses 'input' case-insensitively into a StorageMoverId
// note: this method should only be used for API response data and not user input
func ParseStorageMoverIDInsensitively(input string) (*StorageMoverId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageMoverId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageMoverId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageMoverName, ok = parsed.Parsed["storageMoverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", *parsed)
	}

	return &id, nil
}

// ValidateStorageMoverID checks that 'input' can be parsed as a Storage Mover ID
func ValidateStorageMoverID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageMoverID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Mover ID
func (id StorageMoverId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageMover/storageMovers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Mover ID
func (id StorageMoverId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageMover", "Microsoft.StorageMover", "Microsoft.StorageMover"),
		resourceids.StaticSegment("staticStorageMovers", "storageMovers", "storageMovers"),
		resourceids.UserSpecifiedSegment("storageMoverName", "storageMoverValue"),
	}
}

// String returns a human-readable description of this Storage Mover ID
func (id StorageMoverId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Mover Name: %q", id.StorageMoverName),
	}
	return fmt.Sprintf("Storage Mover (%s)", strings.Join(components, "\n"))
}
