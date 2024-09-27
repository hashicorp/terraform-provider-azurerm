package containerinstance

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ContainerGroupId{})
}

var _ resourceids.ResourceId = &ContainerGroupId{}

// ContainerGroupId is a struct representing the Resource ID for a Container Group
type ContainerGroupId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ContainerGroupName string
}

// NewContainerGroupID returns a new ContainerGroupId struct
func NewContainerGroupID(subscriptionId string, resourceGroupName string, containerGroupName string) ContainerGroupId {
	return ContainerGroupId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ContainerGroupName: containerGroupName,
	}
}

// ParseContainerGroupID parses 'input' into a ContainerGroupId
func ParseContainerGroupID(input string) (*ContainerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContainerGroupIDInsensitively parses 'input' case-insensitively into a ContainerGroupId
// note: this method should only be used for API response data and not user input
func ParseContainerGroupIDInsensitively(input string) (*ContainerGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContainerGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ContainerGroupName, ok = input.Parsed["containerGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerGroupName", input)
	}

	return nil
}

// ValidateContainerGroupID checks that 'input' can be parsed as a Container Group ID
func ValidateContainerGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container Group ID
func (id ContainerGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container Group ID
func (id ContainerGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerInstance", "Microsoft.ContainerInstance", "Microsoft.ContainerInstance"),
		resourceids.StaticSegment("staticContainerGroups", "containerGroups", "containerGroups"),
		resourceids.UserSpecifiedSegment("containerGroupName", "containerGroupName"),
	}
}

// String returns a human-readable description of this Container Group ID
func (id ContainerGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container Group Name: %q", id.ContainerGroupName),
	}
	return fmt.Sprintf("Container Group (%s)", strings.Join(components, "\n"))
}
