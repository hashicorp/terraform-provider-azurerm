package storagetargets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageTargetId{})
}

var _ resourceids.ResourceId = &StorageTargetId{}

// StorageTargetId is a struct representing the Resource ID for a Storage Target
type StorageTargetId struct {
	SubscriptionId    string
	ResourceGroupName string
	CacheName         string
	StorageTargetName string
}

// NewStorageTargetID returns a new StorageTargetId struct
func NewStorageTargetID(subscriptionId string, resourceGroupName string, cacheName string, storageTargetName string) StorageTargetId {
	return StorageTargetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CacheName:         cacheName,
		StorageTargetName: storageTargetName,
	}
}

// ParseStorageTargetID parses 'input' into a StorageTargetId
func ParseStorageTargetID(input string) (*StorageTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageTargetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageTargetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageTargetIDInsensitively parses 'input' case-insensitively into a StorageTargetId
// note: this method should only be used for API response data and not user input
func ParseStorageTargetIDInsensitively(input string) (*StorageTargetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageTargetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageTargetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageTargetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CacheName, ok = input.Parsed["cacheName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cacheName", input)
	}

	if id.StorageTargetName, ok = input.Parsed["storageTargetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageTargetName", input)
	}

	return nil
}

// ValidateStorageTargetID checks that 'input' can be parsed as a Storage Target ID
func ValidateStorageTargetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageTargetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Target ID
func (id StorageTargetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageCache/caches/%s/storageTargets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CacheName, id.StorageTargetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Target ID
func (id StorageTargetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageCache", "Microsoft.StorageCache", "Microsoft.StorageCache"),
		resourceids.StaticSegment("staticCaches", "caches", "caches"),
		resourceids.UserSpecifiedSegment("cacheName", "cacheName"),
		resourceids.StaticSegment("staticStorageTargets", "storageTargets", "storageTargets"),
		resourceids.UserSpecifiedSegment("storageTargetName", "storageTargetName"),
	}
}

// String returns a human-readable description of this Storage Target ID
func (id StorageTargetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cache Name: %q", id.CacheName),
		fmt.Sprintf("Storage Target Name: %q", id.StorageTargetName),
	}
	return fmt.Sprintf("Storage Target (%s)", strings.Join(components, "\n"))
}
