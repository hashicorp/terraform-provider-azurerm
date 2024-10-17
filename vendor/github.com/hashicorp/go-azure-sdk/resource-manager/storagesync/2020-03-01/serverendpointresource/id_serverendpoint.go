package serverendpointresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ServerEndpointId{})
}

var _ resourceids.ResourceId = &ServerEndpointId{}

// ServerEndpointId is a struct representing the Resource ID for a Server Endpoint
type ServerEndpointId struct {
	SubscriptionId         string
	ResourceGroupName      string
	StorageSyncServiceName string
	SyncGroupName          string
	ServerEndpointName     string
}

// NewServerEndpointID returns a new ServerEndpointId struct
func NewServerEndpointID(subscriptionId string, resourceGroupName string, storageSyncServiceName string, syncGroupName string, serverEndpointName string) ServerEndpointId {
	return ServerEndpointId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		StorageSyncServiceName: storageSyncServiceName,
		SyncGroupName:          syncGroupName,
		ServerEndpointName:     serverEndpointName,
	}
}

// ParseServerEndpointID parses 'input' into a ServerEndpointId
func ParseServerEndpointID(input string) (*ServerEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ServerEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ServerEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseServerEndpointIDInsensitively parses 'input' case-insensitively into a ServerEndpointId
// note: this method should only be used for API response data and not user input
func ParseServerEndpointIDInsensitively(input string) (*ServerEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ServerEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ServerEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ServerEndpointId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ServerEndpointName, ok = input.Parsed["serverEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverEndpointName", input)
	}

	return nil
}

// ValidateServerEndpointID checks that 'input' can be parsed as a Server Endpoint ID
func ValidateServerEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseServerEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Server Endpoint ID
func (id ServerEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageSync/storageSyncServices/%s/syncGroups/%s/serverEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName, id.SyncGroupName, id.ServerEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Server Endpoint ID
func (id ServerEndpointId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticServerEndpoints", "serverEndpoints", "serverEndpoints"),
		resourceids.UserSpecifiedSegment("serverEndpointName", "serverEndpointName"),
	}
}

// String returns a human-readable description of this Server Endpoint ID
func (id ServerEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Sync Service Name: %q", id.StorageSyncServiceName),
		fmt.Sprintf("Sync Group Name: %q", id.SyncGroupName),
		fmt.Sprintf("Server Endpoint Name: %q", id.ServerEndpointName),
	}
	return fmt.Sprintf("Server Endpoint (%s)", strings.Join(components, "\n"))
}
