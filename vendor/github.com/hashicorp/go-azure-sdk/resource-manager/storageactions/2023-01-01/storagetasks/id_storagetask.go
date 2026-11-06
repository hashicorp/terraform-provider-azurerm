package storagetasks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageTaskId{})
}

var _ resourceids.ResourceId = &StorageTaskId{}

// StorageTaskId is a struct representing the Resource ID for a Storage Task
type StorageTaskId struct {
	SubscriptionId    string
	ResourceGroupName string
	StorageTaskName   string
}

// NewStorageTaskID returns a new StorageTaskId struct
func NewStorageTaskID(subscriptionId string, resourceGroupName string, storageTaskName string) StorageTaskId {
	return StorageTaskId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StorageTaskName:   storageTaskName,
	}
}

// ParseStorageTaskID parses 'input' into a StorageTaskId
func ParseStorageTaskID(input string) (*StorageTaskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageTaskId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageTaskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageTaskIDInsensitively parses 'input' case-insensitively into a StorageTaskId
// note: this method should only be used for API response data and not user input
func ParseStorageTaskIDInsensitively(input string) (*StorageTaskId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageTaskId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageTaskId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageTaskId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageTaskName, ok = input.Parsed["storageTaskName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageTaskName", input)
	}

	return nil
}

// ValidateStorageTaskID checks that 'input' can be parsed as a Storage Task ID
func ValidateStorageTaskID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageTaskID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Task ID
func (id StorageTaskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageActions/storageTasks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageTaskName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Task ID
func (id StorageTaskId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageActions", "Microsoft.StorageActions", "Microsoft.StorageActions"),
		resourceids.StaticSegment("staticStorageTasks", "storageTasks", "storageTasks"),
		resourceids.UserSpecifiedSegment("storageTaskName", "storageTaskName"),
	}
}

// String returns a human-readable description of this Storage Task ID
func (id StorageTaskId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Task Name: %q", id.StorageTaskName),
	}
	return fmt.Sprintf("Storage Task (%s)", strings.Join(components, "\n"))
}
