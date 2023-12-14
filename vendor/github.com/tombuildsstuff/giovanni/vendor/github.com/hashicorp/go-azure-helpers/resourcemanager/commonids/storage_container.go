// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = StorageContainerId{}

// StorageContainerId is a struct representing the Resource ID for a Storage Container
type StorageContainerId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StorageAccountName string
	ContainerName      string
}

// NewStorageContainerID returns a new NewStorageContainerId struct
func NewStorageContainerID(subscriptionId string, resourceGroupName string, storageAccountName string, containerName string) StorageContainerId {
	return StorageContainerId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		ContainerName:      containerName,
	}
}

// ParseStorageContainerID parses 'input' into a StorageContainerId
func ParseStorageContainerID(input string) (*StorageContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageContainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageContainerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageAccountName, ok = parsed.Parsed["storageAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", *parsed)
	}

	if id.ContainerName, ok = parsed.Parsed["containerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "containerName", *parsed)
	}

	return &id, nil
}

// ParseStorageContainerIDInsensitively parses 'input' case-insensitively into a StorageContainerId
// note: this method should only be used for API response data and not user input
func ParseStorageContainerIDInsensitively(input string) (*StorageContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageContainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageContainerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageAccountName, ok = parsed.Parsed["storageAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", *parsed)
	}

	if id.ContainerName, ok = parsed.Parsed["containerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "containerName", *parsed)
	}

	return &id, nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.ContainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Container ID
func (id StorageContainerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountValue"),
		resourceids.StaticSegment("staticBlobServices", "blobServices", "blobServices"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
		resourceids.StaticSegment("staticContainers", "containers", "containers"),
		resourceids.UserSpecifiedSegment("containerName", "containerValue"),
	}
}

// String returns a human-readable description of this Storage Container ID
func (id StorageContainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Container Name: %q", id.ContainerName),
	}
	return fmt.Sprintf("Storage Container (%s)", strings.Join(components, "\n"))
}
