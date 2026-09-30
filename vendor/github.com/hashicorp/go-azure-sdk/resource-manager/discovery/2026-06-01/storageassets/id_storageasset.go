package storageassets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageAssetId{})
}

var _ resourceids.ResourceId = &StorageAssetId{}

// StorageAssetId is a struct representing the Resource ID for a Storage Asset
type StorageAssetId struct {
	SubscriptionId       string
	ResourceGroupName    string
	StorageContainerName string
	StorageAssetName     string
}

// NewStorageAssetID returns a new StorageAssetId struct
func NewStorageAssetID(subscriptionId string, resourceGroupName string, storageContainerName string, storageAssetName string) StorageAssetId {
	return StorageAssetId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		StorageContainerName: storageContainerName,
		StorageAssetName:     storageAssetName,
	}
}

// ParseStorageAssetID parses 'input' into a StorageAssetId
func ParseStorageAssetID(input string) (*StorageAssetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageAssetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageAssetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageAssetIDInsensitively parses 'input' case-insensitively into a StorageAssetId
// note: this method should only be used for API response data and not user input
func ParseStorageAssetIDInsensitively(input string) (*StorageAssetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageAssetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageAssetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageAssetId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.StorageAssetName, ok = input.Parsed["storageAssetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageAssetName", input)
	}

	return nil
}

// ValidateStorageAssetID checks that 'input' can be parsed as a Storage Asset ID
func ValidateStorageAssetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageAssetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Asset ID
func (id StorageAssetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/storageContainers/%s/storageAssets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageContainerName, id.StorageAssetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Asset ID
func (id StorageAssetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticStorageContainers", "storageContainers", "storageContainers"),
		resourceids.UserSpecifiedSegment("storageContainerName", "storageContainerName"),
		resourceids.StaticSegment("staticStorageAssets", "storageAssets", "storageAssets"),
		resourceids.UserSpecifiedSegment("storageAssetName", "storageAssetName"),
	}
}

// String returns a human-readable description of this Storage Asset ID
func (id StorageAssetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Container Name: %q", id.StorageContainerName),
		fmt.Sprintf("Storage Asset Name: %q", id.StorageAssetName),
	}
	return fmt.Sprintf("Storage Asset (%s)", strings.Join(components, "\n"))
}
