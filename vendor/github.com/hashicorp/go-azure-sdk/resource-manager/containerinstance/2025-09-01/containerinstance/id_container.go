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
	recaser.RegisterResourceId(&ContainerId{})
}

var _ resourceids.ResourceId = &ContainerId{}

// ContainerId is a struct representing the Resource ID for a Container
type ContainerId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ContainerGroupName string
	ContainerName      string
}

// NewContainerID returns a new ContainerId struct
func NewContainerID(subscriptionId string, resourceGroupName string, containerGroupName string, containerName string) ContainerId {
	return ContainerId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ContainerGroupName: containerGroupName,
		ContainerName:      containerName,
	}
}

// ParseContainerID parses 'input' into a ContainerId
func ParseContainerID(input string) (*ContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContainerIDInsensitively parses 'input' case-insensitively into a ContainerId
// note: this method should only be used for API response data and not user input
func ParseContainerIDInsensitively(input string) (*ContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContainerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ContainerName, ok = input.Parsed["containerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerName", input)
	}

	return nil
}

// ValidateContainerID checks that 'input' can be parsed as a Container ID
func ValidateContainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container ID
func (id ContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroups/%s/containers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerGroupName, id.ContainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container ID
func (id ContainerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerInstance", "Microsoft.ContainerInstance", "Microsoft.ContainerInstance"),
		resourceids.StaticSegment("staticContainerGroups", "containerGroups", "containerGroups"),
		resourceids.UserSpecifiedSegment("containerGroupName", "containerGroupName"),
		resourceids.StaticSegment("staticContainers", "containers", "containers"),
		resourceids.UserSpecifiedSegment("containerName", "containerName"),
	}
}

// String returns a human-readable description of this Container ID
func (id ContainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container Group Name: %q", id.ContainerGroupName),
		fmt.Sprintf("Container Name: %q", id.ContainerName),
	}
	return fmt.Sprintf("Container (%s)", strings.Join(components, "\n"))
}
