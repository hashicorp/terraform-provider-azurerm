package storagemovers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageMoverId{})
}

var _ resourceids.ResourceId = &StorageMoverId{}

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
	parser := resourceids.NewParserFromResourceIdType(&StorageMoverId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageMoverId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageMoverIDInsensitively parses 'input' case-insensitively into a StorageMoverId
// note: this method should only be used for API response data and not user input
func ParseStorageMoverIDInsensitively(input string) (*StorageMoverId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageMoverId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageMoverId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageMoverId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageMoverName, ok = input.Parsed["storageMoverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("storageMoverName", "storageMoverName"),
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
