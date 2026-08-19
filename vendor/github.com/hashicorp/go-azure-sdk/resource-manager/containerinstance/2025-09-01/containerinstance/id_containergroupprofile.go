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
	recaser.RegisterResourceId(&ContainerGroupProfileId{})
}

var _ resourceids.ResourceId = &ContainerGroupProfileId{}

// ContainerGroupProfileId is a struct representing the Resource ID for a Container Group Profile
type ContainerGroupProfileId struct {
	SubscriptionId            string
	ResourceGroupName         string
	ContainerGroupProfileName string
}

// NewContainerGroupProfileID returns a new ContainerGroupProfileId struct
func NewContainerGroupProfileID(subscriptionId string, resourceGroupName string, containerGroupProfileName string) ContainerGroupProfileId {
	return ContainerGroupProfileId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		ContainerGroupProfileName: containerGroupProfileName,
	}
}

// ParseContainerGroupProfileID parses 'input' into a ContainerGroupProfileId
func ParseContainerGroupProfileID(input string) (*ContainerGroupProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerGroupProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerGroupProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContainerGroupProfileIDInsensitively parses 'input' case-insensitively into a ContainerGroupProfileId
// note: this method should only be used for API response data and not user input
func ParseContainerGroupProfileIDInsensitively(input string) (*ContainerGroupProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerGroupProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerGroupProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContainerGroupProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ContainerGroupProfileName, ok = input.Parsed["containerGroupProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerGroupProfileName", input)
	}

	return nil
}

// ValidateContainerGroupProfileID checks that 'input' can be parsed as a Container Group Profile ID
func ValidateContainerGroupProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerGroupProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container Group Profile ID
func (id ContainerGroupProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/containerGroupProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerGroupProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container Group Profile ID
func (id ContainerGroupProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerInstance", "Microsoft.ContainerInstance", "Microsoft.ContainerInstance"),
		resourceids.StaticSegment("staticContainerGroupProfiles", "containerGroupProfiles", "containerGroupProfiles"),
		resourceids.UserSpecifiedSegment("containerGroupProfileName", "containerGroupProfileName"),
	}
}

// String returns a human-readable description of this Container Group Profile ID
func (id ContainerGroupProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container Group Profile Name: %q", id.ContainerGroupProfileName),
	}
	return fmt.Sprintf("Container Group Profile (%s)", strings.Join(components, "\n"))
}
