package connectedregistries

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectedRegistryId{})
}

var _ resourceids.ResourceId = &ConnectedRegistryId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ConnectedRegistryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedRegistryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectedRegistryIDInsensitively parses 'input' case-insensitively into a ConnectedRegistryId
// note: this method should only be used for API response data and not user input
func ParseConnectedRegistryIDInsensitively(input string) (*ConnectedRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedRegistryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedRegistryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectedRegistryId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.ConnectedRegistryName, ok = input.Parsed["connectedRegistryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectedRegistryName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticConnectedRegistries", "connectedRegistries", "connectedRegistries"),
		resourceids.UserSpecifiedSegment("connectedRegistryName", "connectedRegistryName"),
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
