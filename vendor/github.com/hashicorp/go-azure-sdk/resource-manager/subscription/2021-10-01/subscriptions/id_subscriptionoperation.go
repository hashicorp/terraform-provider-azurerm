package subscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SubscriptionOperationId{}

// SubscriptionOperationId is a struct representing the Resource ID for a Subscription Operation
type SubscriptionOperationId struct {
	OperationId string
}

// NewSubscriptionOperationID returns a new SubscriptionOperationId struct
func NewSubscriptionOperationID(operationId string) SubscriptionOperationId {
	return SubscriptionOperationId{
		OperationId: operationId,
	}
}

// ParseSubscriptionOperationID parses 'input' into a SubscriptionOperationId
func ParseSubscriptionOperationID(input string) (*SubscriptionOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubscriptionOperationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubscriptionOperationId{}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "operationId", *parsed)
	}

	return &id, nil
}

// ParseSubscriptionOperationIDInsensitively parses 'input' case-insensitively into a SubscriptionOperationId
// note: this method should only be used for API response data and not user input
func ParseSubscriptionOperationIDInsensitively(input string) (*SubscriptionOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubscriptionOperationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubscriptionOperationId{}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "operationId", *parsed)
	}

	return &id, nil
}

// ValidateSubscriptionOperationID checks that 'input' can be parsed as a Subscription Operation ID
func ValidateSubscriptionOperationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubscriptionOperationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subscription Operation ID
func (id SubscriptionOperationId) ID() string {
	fmtString := "/providers/Microsoft.Subscription/subscriptionOperations/%s"
	return fmt.Sprintf(fmtString, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Subscription Operation ID
func (id SubscriptionOperationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSubscription", "Microsoft.Subscription", "Microsoft.Subscription"),
		resourceids.StaticSegment("staticSubscriptionOperations", "subscriptionOperations", "subscriptionOperations"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Subscription Operation ID
func (id SubscriptionOperationId) String() string {
	components := []string{
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Subscription Operation (%s)", strings.Join(components, "\n"))
}
