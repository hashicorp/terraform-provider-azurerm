package edgemodules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VideoAnalyzerId{}

// VideoAnalyzerId is a struct representing the Resource ID for a Video Analyzer
type VideoAnalyzerId struct {
	SubscriptionId    string
	ResourceGroupName string
	VideoAnalyzerName string
}

// NewVideoAnalyzerID returns a new VideoAnalyzerId struct
func NewVideoAnalyzerID(subscriptionId string, resourceGroupName string, videoAnalyzerName string) VideoAnalyzerId {
	return VideoAnalyzerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VideoAnalyzerName: videoAnalyzerName,
	}
}

// ParseVideoAnalyzerID parses 'input' into a VideoAnalyzerId
func ParseVideoAnalyzerID(input string) (*VideoAnalyzerId, error) {
	parser := resourceids.NewParserFromResourceIdType(VideoAnalyzerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VideoAnalyzerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VideoAnalyzerName, ok = parsed.Parsed["videoAnalyzerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "videoAnalyzerName", *parsed)
	}

	return &id, nil
}

// ParseVideoAnalyzerIDInsensitively parses 'input' case-insensitively into a VideoAnalyzerId
// note: this method should only be used for API response data and not user input
func ParseVideoAnalyzerIDInsensitively(input string) (*VideoAnalyzerId, error) {
	parser := resourceids.NewParserFromResourceIdType(VideoAnalyzerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VideoAnalyzerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VideoAnalyzerName, ok = parsed.Parsed["videoAnalyzerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "videoAnalyzerName", *parsed)
	}

	return &id, nil
}

// ValidateVideoAnalyzerID checks that 'input' can be parsed as a Video Analyzer ID
func ValidateVideoAnalyzerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVideoAnalyzerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Video Analyzer ID
func (id VideoAnalyzerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/videoAnalyzers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VideoAnalyzerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Video Analyzer ID
func (id VideoAnalyzerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticVideoAnalyzers", "videoAnalyzers", "videoAnalyzers"),
		resourceids.UserSpecifiedSegment("videoAnalyzerName", "videoAnalyzerValue"),
	}
}

// String returns a human-readable description of this Video Analyzer ID
func (id VideoAnalyzerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Video Analyzer Name: %q", id.VideoAnalyzerName),
	}
	return fmt.Sprintf("Video Analyzer (%s)", strings.Join(components, "\n"))
}
