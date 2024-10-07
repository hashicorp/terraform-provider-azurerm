package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedEventSubscriptionId{})
}

var _ resourceids.ResourceId = &ScopedEventSubscriptionId{}

// ScopedEventSubscriptionId is a struct representing the Resource ID for a Scoped Event Subscription
type ScopedEventSubscriptionId struct {
	Scope                 string
	EventSubscriptionName string
}

// NewScopedEventSubscriptionID returns a new ScopedEventSubscriptionId struct
func NewScopedEventSubscriptionID(scope string, eventSubscriptionName string) ScopedEventSubscriptionId {
	return ScopedEventSubscriptionId{
		Scope:                 scope,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParseScopedEventSubscriptionID parses 'input' into a ScopedEventSubscriptionId
func ParseScopedEventSubscriptionID(input string) (*ScopedEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedEventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedEventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedEventSubscriptionIDInsensitively parses 'input' case-insensitively into a ScopedEventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseScopedEventSubscriptionIDInsensitively(input string) (*ScopedEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedEventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedEventSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedEventSubscriptionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.Scope, ok = input.Parsed["scope"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scope", input)
	}

	if id.EventSubscriptionName, ok = input.Parsed["eventSubscriptionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", input)
	}

	return nil
}

// ValidateScopedEventSubscriptionID checks that 'input' can be parsed as a Scoped Event Subscription ID
func ValidateScopedEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Event Subscription ID
func (id ScopedEventSubscriptionId) ID() string {
	fmtString := "/%s/providers/Microsoft.EventGrid/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Event Subscription ID
func (id ScopedEventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionName"),
	}
}

// String returns a human-readable description of this Scoped Event Subscription ID
func (id ScopedEventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("Scoped Event Subscription (%s)", strings.Join(components, "\n"))
}
