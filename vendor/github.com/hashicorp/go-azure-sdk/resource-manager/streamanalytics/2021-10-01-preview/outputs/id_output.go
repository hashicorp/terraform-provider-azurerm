package outputs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = OutputId{}

// OutputId is a struct representing the Resource ID for a Output
type OutputId struct {
	SubscriptionId    string
	ResourceGroupName string
	StreamingJobName  string
	OutputName        string
}

// NewOutputID returns a new OutputId struct
func NewOutputID(subscriptionId string, resourceGroupName string, streamingJobName string, outputName string) OutputId {
	return OutputId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StreamingJobName:  streamingJobName,
		OutputName:        outputName,
	}
}

// ParseOutputID parses 'input' into a OutputId
func ParseOutputID(input string) (*OutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(OutputId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OutputId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamingJobName", *parsed)
	}

	if id.OutputName, ok = parsed.Parsed["outputName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "outputName", *parsed)
	}

	return &id, nil
}

// ParseOutputIDInsensitively parses 'input' case-insensitively into a OutputId
// note: this method should only be used for API response data and not user input
func ParseOutputIDInsensitively(input string) (*OutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(OutputId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OutputId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "streamingJobName", *parsed)
	}

	if id.OutputName, ok = parsed.Parsed["outputName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "outputName", *parsed)
	}

	return &id, nil
}

// ValidateOutputID checks that 'input' can be parsed as a Output ID
func ValidateOutputID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOutputID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Output ID
func (id OutputId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingJobs/%s/outputs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName, id.OutputName)
}

// Segments returns a slice of Resource ID Segments which comprise this Output ID
func (id OutputId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStreamAnalytics", "Microsoft.StreamAnalytics", "Microsoft.StreamAnalytics"),
		resourceids.StaticSegment("staticStreamingJobs", "streamingJobs", "streamingJobs"),
		resourceids.UserSpecifiedSegment("streamingJobName", "streamingJobValue"),
		resourceids.StaticSegment("staticOutputs", "outputs", "outputs"),
		resourceids.UserSpecifiedSegment("outputName", "outputValue"),
	}
}

// String returns a human-readable description of this Output ID
func (id OutputId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Streaming Job Name: %q", id.StreamingJobName),
		fmt.Sprintf("Output Name: %q", id.OutputName),
	}
	return fmt.Sprintf("Output (%s)", strings.Join(components, "\n"))
}
