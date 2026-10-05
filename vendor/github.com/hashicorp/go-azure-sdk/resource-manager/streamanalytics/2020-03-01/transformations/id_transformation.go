package transformations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TransformationId{})
}

var _ resourceids.ResourceId = &TransformationId{}

// TransformationId is a struct representing the Resource ID for a Transformation
type TransformationId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StreamingJobName   string
	TransformationName string
}

// NewTransformationID returns a new TransformationId struct
func NewTransformationID(subscriptionId string, resourceGroupName string, streamingJobName string, transformationName string) TransformationId {
	return TransformationId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StreamingJobName:   streamingJobName,
		TransformationName: transformationName,
	}
}

// ParseTransformationID parses 'input' into a TransformationId
func ParseTransformationID(input string) (*TransformationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TransformationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TransformationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTransformationIDInsensitively parses 'input' case-insensitively into a TransformationId
// note: this method should only be used for API response data and not user input
func ParseTransformationIDInsensitively(input string) (*TransformationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TransformationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TransformationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TransformationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StreamingJobName, ok = input.Parsed["streamingJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "streamingJobName", input)
	}

	if id.TransformationName, ok = input.Parsed["transformationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "transformationName", input)
	}

	return nil
}

// ValidateTransformationID checks that 'input' can be parsed as a Transformation ID
func ValidateTransformationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTransformationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Transformation ID
func (id TransformationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingJobs/%s/transformations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName, id.TransformationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Transformation ID
func (id TransformationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStreamAnalytics", "Microsoft.StreamAnalytics", "Microsoft.StreamAnalytics"),
		resourceids.StaticSegment("staticStreamingJobs", "streamingJobs", "streamingJobs"),
		resourceids.UserSpecifiedSegment("streamingJobName", "streamingJobName"),
		resourceids.StaticSegment("staticTransformations", "transformations", "transformations"),
		resourceids.UserSpecifiedSegment("transformationName", "transformationName"),
	}
}

// String returns a human-readable description of this Transformation ID
func (id TransformationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Streaming Job Name: %q", id.StreamingJobName),
		fmt.Sprintf("Transformation Name: %q", id.TransformationName),
	}
	return fmt.Sprintf("Transformation (%s)", strings.Join(components, "\n"))
}
