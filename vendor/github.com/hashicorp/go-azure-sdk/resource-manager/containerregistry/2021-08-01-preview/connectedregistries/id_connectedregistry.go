package connectedregistries

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ConnectedRegistryId{}

// ConnectedRegistryId is a struct representing the Resource ID for a Connected Registry
type ConnectedRegistryId struct {
	SubscriptionId        string
	ResourceGroupName     string
	RegistryName          string
	ConnectedRegistryName string
}

// NewConnectedRegistryID returns a new ConnectedRegistryId struct
func NewConnectedRegistryID(subscriptionId string, resourceGroupName string, registryName string, connectedRegistryName string) ConnectedRegistryId {
	return ConnectedRegistryId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		RegistryName:          registryName,
		ConnectedRegistryName: connectedRegistryName,
	}
}

// ParseConnectedRegistryID parses 'input' into a ConnectedRegistryId
func ParseConnectedRegistryID(input string) (*ConnectedRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectedRegistryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectedRegistryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ConnectedRegistryName, ok = parsed.Parsed["connectedRegistryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "connectedRegistryName", *parsed)
	}

	return &id, nil
}

// ParseConnectedRegistryIDInsensitively parses 'input' case-insensitively into a ConnectedRegistryId
// note: this method should only be used for API response data and not user input
func ParseConnectedRegistryIDInsensitively(input string) (*ConnectedRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectedRegistryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectedRegistryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ConnectedRegistryName, ok = parsed.Parsed["connectedRegistryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "connectedRegistryName", *parsed)
	}

	return &id, nil
}

// ValidateConnectedRegistryID checks that 'input' can be parsed as a Connected Registry ID
func ValidateConnectedRegistryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectedRegistryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connected Registry ID
func (id ConnectedRegistryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/connectedRegistries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ConnectedRegistryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connected Registry ID
func (id ConnectedRegistryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticConnectedRegistries", "connectedRegistries", "connectedRegistries"),
		resourceids.UserSpecifiedSegment("connectedRegistryName", "connectedRegistryValue"),
	}
}

// String returns a human-readable description of this Connected Registry ID
func (id ConnectedRegistryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Connected Registry Name: %q", id.ConnectedRegistryName),
	}
	return fmt.Sprintf("Connected Registry (%s)", strings.Join(components, "\n"))
}
