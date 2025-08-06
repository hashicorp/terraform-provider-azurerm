package operationprogress

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OperationProgressId{})
}

var _ resourceids.ResourceId = &OperationProgressId{}

// OperationProgressId is a struct representing the Resource ID for a Operation Progress
type OperationProgressId struct {
	SubscriptionId string
	LocationName   string
	OperationId    string
}

// NewOperationProgressID returns a new OperationProgressId struct
func NewOperationProgressID(subscriptionId string, locationName string, operationId string) OperationProgressId {
	return OperationProgressId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		OperationId:    operationId,
	}
}

// ParseOperationProgressID parses 'input' into a OperationProgressId
func ParseOperationProgressID(input string) (*OperationProgressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationProgressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationProgressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOperationProgressIDInsensitively parses 'input' case-insensitively into a OperationProgressId
// note: this method should only be used for API response data and not user input
func ParseOperationProgressIDInsensitively(input string) (*OperationProgressId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OperationProgressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OperationProgressId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OperationProgressId) FromParseResult(input resourceids.ParseResult) error {
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

// ValidateOperationProgressID checks that 'input' can be parsed as a Operation Progress ID
func ValidateOperationProgressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationProgressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Progress ID
func (id OperationProgressId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DBforMySQL/locations/%s/operationProgress/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Progress ID
func (id OperationProgressId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMySQL", "Microsoft.DBforMySQL", "Microsoft.DBforMySQL"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticOperationProgress", "operationProgress", "operationProgress"),
		resourceids.UserSpecifiedSegment("operationId", "operationId"),
	}
}

// String returns a human-readable description of this Operation Progress ID
func (id OperationProgressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Operation Progress (%s)", strings.Join(components, "\n"))
}
