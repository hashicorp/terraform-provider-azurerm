package managedenvironments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedEnvironmentDetectorId{})
}

var _ resourceids.ResourceId = &ManagedEnvironmentDetectorId{}

// ManagedEnvironmentDetectorId is a struct representing the Resource ID for a Managed Environment Detector
type ManagedEnvironmentDetectorId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	DetectorName           string
}

// NewManagedEnvironmentDetectorID returns a new ManagedEnvironmentDetectorId struct
func NewManagedEnvironmentDetectorID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, detectorName string) ManagedEnvironmentDetectorId {
	return ManagedEnvironmentDetectorId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		DetectorName:           detectorName,
	}
}

// ParseManagedEnvironmentDetectorID parses 'input' into a ManagedEnvironmentDetectorId
func ParseManagedEnvironmentDetectorID(input string) (*ManagedEnvironmentDetectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedEnvironmentDetectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedEnvironmentDetectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedEnvironmentDetectorIDInsensitively parses 'input' case-insensitively into a ManagedEnvironmentDetectorId
// note: this method should only be used for API response data and not user input
func ParseManagedEnvironmentDetectorIDInsensitively(input string) (*ManagedEnvironmentDetectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedEnvironmentDetectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedEnvironmentDetectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedEnvironmentDetectorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedEnvironmentName, ok = input.Parsed["managedEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedEnvironmentName", input)
	}

	if id.DetectorName, ok = input.Parsed["detectorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "detectorName", input)
	}

	return nil
}

// ValidateManagedEnvironmentDetectorID checks that 'input' can be parsed as a Managed Environment Detector ID
func ValidateManagedEnvironmentDetectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedEnvironmentDetectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Environment Detector ID
func (id ManagedEnvironmentDetectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/detectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.DetectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Environment Detector ID
func (id ManagedEnvironmentDetectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentName"),
		resourceids.StaticSegment("staticDetectors", "detectors", "detectors"),
		resourceids.UserSpecifiedSegment("detectorName", "detectorName"),
	}
}

// String returns a human-readable description of this Managed Environment Detector ID
func (id ManagedEnvironmentDetectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Detector Name: %q", id.DetectorName),
	}
	return fmt.Sprintf("Managed Environment Detector (%s)", strings.Join(components, "\n"))
}
