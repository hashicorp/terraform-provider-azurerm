package accounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = MediaServicesOperationResultId{}

// MediaServicesOperationResultId is a struct representing the Resource ID for a Media Services Operation Result
type MediaServicesOperationResultId struct {
	SubscriptionId string
	LocationName   string
	OperationId    string
}

// NewMediaServicesOperationResultID returns a new MediaServicesOperationResultId struct
func NewMediaServicesOperationResultID(subscriptionId string, locationName string, operationId string) MediaServicesOperationResultId {
	return MediaServicesOperationResultId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		OperationId:    operationId,
	}
}

// ParseMediaServicesOperationResultID parses 'input' into a MediaServicesOperationResultId
func ParseMediaServicesOperationResultID(input string) (*MediaServicesOperationResultId, error) {
	parser := resourceids.NewParserFromResourceIdType(MediaServicesOperationResultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MediaServicesOperationResultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "operationId", *parsed)
	}

	return &id, nil
}

// ParseMediaServicesOperationResultIDInsensitively parses 'input' case-insensitively into a MediaServicesOperationResultId
// note: this method should only be used for API response data and not user input
func ParseMediaServicesOperationResultIDInsensitively(input string) (*MediaServicesOperationResultId, error) {
	parser := resourceids.NewParserFromResourceIdType(MediaServicesOperationResultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MediaServicesOperationResultId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "operationId", *parsed)
	}

	return &id, nil
}

// ValidateMediaServicesOperationResultID checks that 'input' can be parsed as a Media Services Operation Result ID
func ValidateMediaServicesOperationResultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMediaServicesOperationResultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Media Services Operation Result ID
func (id MediaServicesOperationResultId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Media/locations/%s/mediaServicesOperationResults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Media Services Operation Result ID
func (id MediaServicesOperationResultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticMediaServicesOperationResults", "mediaServicesOperationResults", "mediaServicesOperationResults"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Media Services Operation Result ID
func (id MediaServicesOperationResultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Media Services Operation Result (%s)", strings.Join(components, "\n"))
}
