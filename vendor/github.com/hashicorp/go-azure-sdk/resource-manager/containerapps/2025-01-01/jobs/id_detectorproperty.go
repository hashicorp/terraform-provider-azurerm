package jobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DetectorPropertyId{})
}

var _ resourceids.ResourceId = &DetectorPropertyId{}

// DetectorPropertyId is a struct representing the Resource ID for a Detector Property
type DetectorPropertyId struct {
	SubscriptionId       string
	ResourceGroupName    string
	JobName              string
	DetectorPropertyName string
}

// NewDetectorPropertyID returns a new DetectorPropertyId struct
func NewDetectorPropertyID(subscriptionId string, resourceGroupName string, jobName string, detectorPropertyName string) DetectorPropertyId {
	return DetectorPropertyId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		JobName:              jobName,
		DetectorPropertyName: detectorPropertyName,
	}
}

// ParseDetectorPropertyID parses 'input' into a DetectorPropertyId
func ParseDetectorPropertyID(input string) (*DetectorPropertyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DetectorPropertyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DetectorPropertyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDetectorPropertyIDInsensitively parses 'input' case-insensitively into a DetectorPropertyId
// note: this method should only be used for API response data and not user input
func ParseDetectorPropertyIDInsensitively(input string) (*DetectorPropertyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DetectorPropertyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DetectorPropertyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DetectorPropertyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.JobName, ok = input.Parsed["jobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobName", input)
	}

	if id.DetectorPropertyName, ok = input.Parsed["detectorPropertyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "detectorPropertyName", input)
	}

	return nil
}

// ValidateDetectorPropertyID checks that 'input' can be parsed as a Detector Property ID
func ValidateDetectorPropertyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDetectorPropertyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Detector Property ID
func (id DetectorPropertyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/jobs/%s/detectorProperties/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.JobName, id.DetectorPropertyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Detector Property ID
func (id DetectorPropertyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobName"),
		resourceids.StaticSegment("staticDetectorProperties", "detectorProperties", "detectorProperties"),
		resourceids.UserSpecifiedSegment("detectorPropertyName", "detectorPropertyName"),
	}
}

// String returns a human-readable description of this Detector Property ID
func (id DetectorPropertyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Job Name: %q", id.JobName),
		fmt.Sprintf("Detector Property Name: %q", id.DetectorPropertyName),
	}
	return fmt.Sprintf("Detector Property (%s)", strings.Join(components, "\n"))
}
