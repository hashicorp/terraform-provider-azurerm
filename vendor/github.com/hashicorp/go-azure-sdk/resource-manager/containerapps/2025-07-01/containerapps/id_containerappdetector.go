package containerapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ContainerAppDetectorId{})
}

var _ resourceids.ResourceId = &ContainerAppDetectorId{}

// ContainerAppDetectorId is a struct representing the Resource ID for a Container App Detector
type ContainerAppDetectorId struct {
	SubscriptionId    string
	ResourceGroupName string
	ContainerAppName  string
	DetectorName      string
}

// NewContainerAppDetectorID returns a new ContainerAppDetectorId struct
func NewContainerAppDetectorID(subscriptionId string, resourceGroupName string, containerAppName string, detectorName string) ContainerAppDetectorId {
	return ContainerAppDetectorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ContainerAppName:  containerAppName,
		DetectorName:      detectorName,
	}
}

// ParseContainerAppDetectorID parses 'input' into a ContainerAppDetectorId
func ParseContainerAppDetectorID(input string) (*ContainerAppDetectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerAppDetectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerAppDetectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseContainerAppDetectorIDInsensitively parses 'input' case-insensitively into a ContainerAppDetectorId
// note: this method should only be used for API response data and not user input
func ParseContainerAppDetectorIDInsensitively(input string) (*ContainerAppDetectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerAppDetectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerAppDetectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContainerAppDetectorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ContainerAppName, ok = input.Parsed["containerAppName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerAppName", input)
	}

	if id.DetectorName, ok = input.Parsed["detectorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "detectorName", input)
	}

	return nil
}

// ValidateContainerAppDetectorID checks that 'input' can be parsed as a Container App Detector ID
func ValidateContainerAppDetectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerAppDetectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container App Detector ID
func (id ContainerAppDetectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/containerApps/%s/detectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName, id.DetectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container App Detector ID
func (id ContainerAppDetectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticContainerApps", "containerApps", "containerApps"),
		resourceids.UserSpecifiedSegment("containerAppName", "containerAppName"),
		resourceids.StaticSegment("staticDetectors", "detectors", "detectors"),
		resourceids.UserSpecifiedSegment("detectorName", "detectorName"),
	}
}

// String returns a human-readable description of this Container App Detector ID
func (id ContainerAppDetectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container App Name: %q", id.ContainerAppName),
		fmt.Sprintf("Detector Name: %q", id.DetectorName),
	}
	return fmt.Sprintf("Container App Detector (%s)", strings.Join(components, "\n"))
}
