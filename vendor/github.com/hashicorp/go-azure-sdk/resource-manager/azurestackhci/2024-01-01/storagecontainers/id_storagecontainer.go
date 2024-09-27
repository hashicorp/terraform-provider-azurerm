package storagecontainers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageContainerId{})
}

var _ resourceids.ResourceId = &StorageContainerId{}

// StorageContainerId is a struct representing the Resource ID for a Storage Container
type StorageContainerId struct {
	SubscriptionId       string
	ResourceGroupName    string
	StorageContainerName string
}

// NewStorageContainerID returns a new StorageContainerId struct
func NewStorageContainerID(subscriptionId string, resourceGroupName string, storageContainerName string) StorageContainerId {
	return StorageContainerId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		StorageContainerName: storageContainerName,
	}
}

// ParseStorageContainerID parses 'input' into a StorageContainerId
func ParseStorageContainerID(input string) (*StorageContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageContainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageContainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageContainerIDInsensitively parses 'input' case-insensitively into a StorageContainerId
// note: this method should only be used for API response data and not user input
func ParseStorageContainerIDInsensitively(input string) (*StorageContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageContainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageContainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageContainerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageContainerName, ok = input.Parsed["storageContainerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageContainerName", input)
	}

	return nil
}

// ValidateStorageContainerID checks that 'input' can be parsed as a Storage Container ID
func ValidateStorageContainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageContainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Container ID
func (id StorageContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/storageContainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageContainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Container ID
func (id StorageContainerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticStorageContainers", "storageContainers", "storageContainers"),
		resourceids.UserSpecifiedSegment("storageContainerName", "storageContainerName"),
	}
}

// String returns a human-readable description of this Storage Container ID
func (id StorageContainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Container Name: %q", id.StorageContainerName),
	}
	return fmt.Sprintf("Storage Container (%s)", strings.Join(components, "\n"))
}
