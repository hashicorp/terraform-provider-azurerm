package redis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AsyncOperationId{})
}

var _ resourceids.ResourceId = &AsyncOperationId{}

// AsyncOperationId is a struct representing the Resource ID for a Async Operation
type AsyncOperationId struct {
	SubscriptionId string
	LocationName   string
	OperationId    string
}

// NewAsyncOperationID returns a new AsyncOperationId struct
func NewAsyncOperationID(subscriptionId string, locationName string, operationId string) AsyncOperationId {
	return AsyncOperationId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		OperationId:    operationId,
	}
}

// ParseAsyncOperationID parses 'input' into a AsyncOperationId
func ParseAsyncOperationID(input string) (*AsyncOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AsyncOperationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AsyncOperationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAsyncOperationIDInsensitively parses 'input' case-insensitively into a AsyncOperationId
// note: this method should only be used for API response data and not user input
func ParseAsyncOperationIDInsensitively(input string) (*AsyncOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AsyncOperationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AsyncOperationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AsyncOperationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.OperationId, ok = input.Parsed["operationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationId", input)
	}

	return nil
}

// ValidateAsyncOperationID checks that 'input' can be parsed as a Async Operation ID
func ValidateAsyncOperationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAsyncOperationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Async Operation ID
func (id AsyncOperationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Cache/locations/%s/asyncOperations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Async Operation ID
func (id AsyncOperationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticAsyncOperations", "asyncOperations", "asyncOperations"),
		resourceids.UserSpecifiedSegment("operationId", "operationId"),
	}
}

// String returns a human-readable description of this Async Operation ID
func (id AsyncOperationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Async Operation (%s)", strings.Join(components, "\n"))
}
