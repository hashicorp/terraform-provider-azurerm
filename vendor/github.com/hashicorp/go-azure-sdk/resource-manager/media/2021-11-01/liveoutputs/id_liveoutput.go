package liveoutputs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LiveOutputId{})
}

var _ resourceids.ResourceId = &LiveOutputId{}

// LiveOutputId is a struct representing the Resource ID for a Live Output
type LiveOutputId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	LiveEventName     string
	LiveOutputName    string
}

// NewLiveOutputID returns a new LiveOutputId struct
func NewLiveOutputID(subscriptionId string, resourceGroupName string, mediaServiceName string, liveEventName string, liveOutputName string) LiveOutputId {
	return LiveOutputId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		LiveEventName:     liveEventName,
		LiveOutputName:    liveOutputName,
	}
}

// ParseLiveOutputID parses 'input' into a LiveOutputId
func ParseLiveOutputID(input string) (*LiveOutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LiveOutputId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LiveOutputId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLiveOutputIDInsensitively parses 'input' case-insensitively into a LiveOutputId
// note: this method should only be used for API response data and not user input
func ParseLiveOutputIDInsensitively(input string) (*LiveOutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LiveOutputId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LiveOutputId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LiveOutputId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MediaServiceName, ok = input.Parsed["mediaServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", input)
	}

	if id.LiveEventName, ok = input.Parsed["liveEventName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "liveEventName", input)
	}

	if id.LiveOutputName, ok = input.Parsed["liveOutputName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "liveOutputName", input)
	}

	return nil
}

// ValidateLiveOutputID checks that 'input' can be parsed as a Live Output ID
func ValidateLiveOutputID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLiveOutputID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Live Output ID
func (id LiveOutputId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/liveEvents/%s/liveOutputs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.LiveEventName, id.LiveOutputName)
}

// Segments returns a slice of Resource ID Segments which comprise this Live Output ID
func (id LiveOutputId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticLiveEvents", "liveEvents", "liveEvents"),
		resourceids.UserSpecifiedSegment("liveEventName", "liveEventValue"),
		resourceids.StaticSegment("staticLiveOutputs", "liveOutputs", "liveOutputs"),
		resourceids.UserSpecifiedSegment("liveOutputName", "liveOutputValue"),
	}
}

// String returns a human-readable description of this Live Output ID
func (id LiveOutputId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Live Event Name: %q", id.LiveEventName),
		fmt.Sprintf("Live Output Name: %q", id.LiveOutputName),
	}
	return fmt.Sprintf("Live Output (%s)", strings.Join(components, "\n"))
}
