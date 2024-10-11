package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ContainerRegistryId{})
}

var _ resourceids.ResourceId = &ContainerRegistryId{}

// ContainerRegistryId is a struct representing the Resource ID for a Container Registry
type ContainerRegistryId struct {
	SubscriptionId        string
	ResourceGroupName     string
	SpringName            string
	ContainerRegistryName string
}

// NewContainerRegistryID returns a new ContainerRegistryId struct
func NewContainerRegistryID(subscriptionId string, resourceGroupName string, springName string, containerRegistryName string) ContainerRegistryId {
	return ContainerRegistryId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		SpringName:            springName,
		ContainerRegistryName: containerRegistryName,
	}
}

// ParseContainerRegistryID parses 'input' into a ContainerRegistryId
func ParseContainerRegistryID(input string) (*ContainerRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerRegistryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerRegistryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContainerRegistryIDInsensitively parses 'input' case-insensitively into a ContainerRegistryId
// note: this method should only be used for API response data and not user input
func ParseContainerRegistryIDInsensitively(input string) (*ContainerRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerRegistryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerRegistryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContainerRegistryId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.ContainerRegistryName, ok = input.Parsed["containerRegistryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerRegistryName", input)
	}

	return nil
}

// ValidateContainerRegistryID checks that 'input' can be parsed as a Container Registry ID
func ValidateContainerRegistryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerRegistryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container Registry ID
func (id ContainerRegistryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/containerRegistries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ContainerRegistryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container Registry ID
func (id ContainerRegistryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticContainerRegistries", "containerRegistries", "containerRegistries"),
		resourceids.UserSpecifiedSegment("containerRegistryName", "containerRegistryName"),
	}
}

// String returns a human-readable description of this Container Registry ID
func (id ContainerRegistryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Container Registry Name: %q", id.ContainerRegistryName),
	}
	return fmt.Sprintf("Container Registry (%s)", strings.Join(components, "\n"))
}
