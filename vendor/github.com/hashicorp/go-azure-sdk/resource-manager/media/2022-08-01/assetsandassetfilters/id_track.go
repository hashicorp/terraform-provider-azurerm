package assetsandassetfilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TrackId{}

// TrackId is a struct representing the Resource ID for a Track
type TrackId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	AssetName         string
	TrackName         string
}

// NewTrackID returns a new TrackId struct
func NewTrackID(subscriptionId string, resourceGroupName string, mediaServiceName string, assetName string, trackName string) TrackId {
	return TrackId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		AssetName:         assetName,
		TrackName:         trackName,
	}
}

// ParseTrackID parses 'input' into a TrackId
func ParseTrackID(input string) (*TrackId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrackId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.AssetName, ok = parsed.Parsed["assetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assetName", *parsed)
	}

	if id.TrackName, ok = parsed.Parsed["trackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trackName", *parsed)
	}

	return &id, nil
}

// ParseTrackIDInsensitively parses 'input' case-insensitively into a TrackId
// note: this method should only be used for API response data and not user input
func ParseTrackIDInsensitively(input string) (*TrackId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrackId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.AssetName, ok = parsed.Parsed["assetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assetName", *parsed)
	}

	if id.TrackName, ok = parsed.Parsed["trackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trackName", *parsed)
	}

	return &id, nil
}

// ValidateTrackID checks that 'input' can be parsed as a Track ID
func ValidateTrackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTrackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Track ID
func (id TrackId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/assets/%s/tracks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.AssetName, id.TrackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Track ID
func (id TrackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticAssets", "assets", "assets"),
		resourceids.UserSpecifiedSegment("assetName", "assetValue"),
		resourceids.StaticSegment("staticTracks", "tracks", "tracks"),
		resourceids.UserSpecifiedSegment("trackName", "trackValue"),
	}
}

// String returns a human-readable description of this Track ID
func (id TrackId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Asset Name: %q", id.AssetName),
		fmt.Sprintf("Track Name: %q", id.TrackName),
	}
	return fmt.Sprintf("Track (%s)", strings.Join(components, "\n"))
}
