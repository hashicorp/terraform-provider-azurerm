package cloudendpointresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CloudEndpointId{})
}

var _ resourceids.ResourceId = &CloudEndpointId{}

// CloudEndpointId is a struct representing the Resource ID for a Cloud Endpoint
type CloudEndpointId struct {
	SubscriptionId         string
	ResourceGroupName      string
	StorageSyncServiceName string
	SyncGroupName          string
	CloudEndpointName      string
}

// NewCloudEndpointID returns a new CloudEndpointId struct
func NewCloudEndpointID(subscriptionId string, resourceGroupName string, storageSyncServiceName string, syncGroupName string, cloudEndpointName string) CloudEndpointId {
	return CloudEndpointId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		StorageSyncServiceName: storageSyncServiceName,
		SyncGroupName:          syncGroupName,
		CloudEndpointName:      cloudEndpointName,
	}
}

// ParseCloudEndpointID parses 'input' into a CloudEndpointId
func ParseCloudEndpointID(input string) (*CloudEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudEndpointIDInsensitively parses 'input' case-insensitively into a CloudEndpointId
// note: this method should only be used for API response data and not user input
func ParseCloudEndpointIDInsensitively(input string) (*CloudEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudEndpointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageSyncServiceName, ok = input.Parsed["storageSyncServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageSyncServiceName", input)
	}

	if id.SyncGroupName, ok = input.Parsed["syncGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "syncGroupName", input)
	}

	if id.CloudEndpointName, ok = input.Parsed["cloudEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudEndpointName", input)
	}

	return nil
}

// ValidateCloudEndpointID checks that 'input' can be parsed as a Cloud Endpoint ID
func ValidateCloudEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud Endpoint ID
func (id CloudEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageSync/storageSyncServices/%s/syncGroups/%s/cloudEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName, id.SyncGroupName, id.CloudEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Endpoint ID
func (id CloudEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageSync", "Microsoft.StorageSync", "Microsoft.StorageSync"),
		resourceids.StaticSegment("staticStorageSyncServices", "storageSyncServices", "storageSyncServices"),
		resourceids.UserSpecifiedSegment("storageSyncServiceName", "storageSyncServiceName"),
		resourceids.StaticSegment("staticSyncGroups", "syncGroups", "syncGroups"),
		resourceids.UserSpecifiedSegment("syncGroupName", "syncGroupName"),
		resourceids.StaticSegment("staticCloudEndpoints", "cloudEndpoints", "cloudEndpoints"),
		resourceids.UserSpecifiedSegment("cloudEndpointName", "cloudEndpointName"),
	}
}

// String returns a human-readable description of this Cloud Endpoint ID
func (id CloudEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Sync Service Name: %q", id.StorageSyncServiceName),
		fmt.Sprintf("Sync Group Name: %q", id.SyncGroupName),
		fmt.Sprintf("Cloud Endpoint Name: %q", id.CloudEndpointName),
	}
	return fmt.Sprintf("Cloud Endpoint (%s)", strings.Join(components, "\n"))
}
