package endpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = EndpointId{}

// EndpointId is a struct representing the Resource ID for a Endpoint
type EndpointId struct {
	SubscriptionId    string
	ResourceGroupName string
	StorageMoverName  string
	EndpointName      string
}

// NewEndpointID returns a new EndpointId struct
func NewEndpointID(subscriptionId string, resourceGroupName string, storageMoverName string, endpointName string) EndpointId {
	return EndpointId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StorageMoverName:  storageMoverName,
		EndpointName:      endpointName,
	}
}

// ParseEndpointID parses 'input' into a EndpointId
func ParseEndpointID(input string) (*EndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(EndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageMoverName, ok = parsed.Parsed["storageMoverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", *parsed)
	}

	if id.EndpointName, ok = parsed.Parsed["endpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "endpointName", *parsed)
	}

	return &id, nil
}

// ParseEndpointIDInsensitively parses 'input' case-insensitively into a EndpointId
// note: this method should only be used for API response data and not user input
func ParseEndpointIDInsensitively(input string) (*EndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(EndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StorageMoverName, ok = parsed.Parsed["storageMoverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", *parsed)
	}

	if id.EndpointName, ok = parsed.Parsed["endpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "endpointName", *parsed)
	}

	return &id, nil
}

// ValidateEndpointID checks that 'input' can be parsed as a Endpoint ID
func ValidateEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Endpoint ID
func (id EndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageMover/storageMovers/%s/endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName, id.EndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Endpoint ID
func (id EndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageMover", "Microsoft.StorageMover", "Microsoft.StorageMover"),
		resourceids.StaticSegment("staticStorageMovers", "storageMovers", "storageMovers"),
		resourceids.UserSpecifiedSegment("storageMoverName", "storageMoverValue"),
		resourceids.StaticSegment("staticEndpoints", "endpoints", "endpoints"),
		resourceids.UserSpecifiedSegment("endpointName", "endpointValue"),
	}
}

// String returns a human-readable description of this Endpoint ID
func (id EndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Mover Name: %q", id.StorageMoverName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
	}
	return fmt.Sprintf("Endpoint (%s)", strings.Join(components, "\n"))
}
