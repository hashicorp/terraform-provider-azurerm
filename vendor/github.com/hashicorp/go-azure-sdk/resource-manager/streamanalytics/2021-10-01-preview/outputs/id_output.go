package outputs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OutputId{})
}

var _ resourceids.ResourceId = &OutputId{}

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
	parser := resourceids.NewParserFromResourceIdType(&OutputId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OutputId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOutputIDInsensitively parses 'input' case-insensitively into a OutputId
// note: this method should only be used for API response data and not user input
func ParseOutputIDInsensitively(input string) (*OutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OutputId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OutputId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OutputId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.OutputName, ok = input.Parsed["outputName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "outputName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("streamingJobName", "streamingJobName"),
		resourceids.StaticSegment("staticOutputs", "outputs", "outputs"),
		resourceids.UserSpecifiedSegment("outputName", "outputName"),
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
