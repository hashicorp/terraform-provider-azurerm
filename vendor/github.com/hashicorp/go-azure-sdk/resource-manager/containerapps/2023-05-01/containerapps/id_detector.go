package containerapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DetectorId{}

// DetectorId is a struct representing the Resource ID for a Detector
type DetectorId struct {
	SubscriptionId    string
	ResourceGroupName string
	ContainerAppName  string
	DetectorName      string
}

// NewDetectorID returns a new DetectorId struct
func NewDetectorID(subscriptionId string, resourceGroupName string, containerAppName string, detectorName string) DetectorId {
	return DetectorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ContainerAppName:  containerAppName,
		DetectorName:      detectorName,
	}
}

// ParseDetectorID parses 'input' into a DetectorId
func ParseDetectorID(input string) (*DetectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(DetectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DetectorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ContainerAppName, ok = parsed.Parsed["containerAppName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "containerAppName", *parsed)
	}

	if id.DetectorName, ok = parsed.Parsed["detectorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "detectorName", *parsed)
	}

	return &id, nil
}

// ParseDetectorIDInsensitively parses 'input' case-insensitively into a DetectorId
// note: this method should only be used for API response data and not user input
func ParseDetectorIDInsensitively(input string) (*DetectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(DetectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DetectorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ContainerAppName, ok = parsed.Parsed["containerAppName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "containerAppName", *parsed)
	}

	if id.DetectorName, ok = parsed.Parsed["detectorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "detectorName", *parsed)
	}

	return &id, nil
}

// ValidateDetectorID checks that 'input' can be parsed as a Detector ID
func ValidateDetectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDetectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Detector ID
func (id DetectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/containerApps/%s/detectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName, id.DetectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Detector ID
func (id DetectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticContainerApps", "containerApps", "containerApps"),
		resourceids.UserSpecifiedSegment("containerAppName", "containerAppValue"),
		resourceids.StaticSegment("staticDetectors", "detectors", "detectors"),
		resourceids.UserSpecifiedSegment("detectorName", "detectorValue"),
	}
}

// String returns a human-readable description of this Detector ID
func (id DetectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container App Name: %q", id.ContainerAppName),
		fmt.Sprintf("Detector Name: %q", id.DetectorName),
	}
	return fmt.Sprintf("Detector (%s)", strings.Join(components, "\n"))
}
