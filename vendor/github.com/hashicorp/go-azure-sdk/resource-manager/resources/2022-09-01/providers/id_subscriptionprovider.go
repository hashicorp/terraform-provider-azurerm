package providers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SubscriptionProviderId{}

// SubscriptionProviderId is a struct representing the Resource ID for a Subscription Provider
type SubscriptionProviderId struct {
	SubscriptionId string
	ProviderName   string
}

// NewSubscriptionProviderID returns a new SubscriptionProviderId struct
func NewSubscriptionProviderID(subscriptionId string, providerName string) SubscriptionProviderId {
	return SubscriptionProviderId{
		SubscriptionId: subscriptionId,
		ProviderName:   providerName,
	}
}

// ParseSubscriptionProviderID parses 'input' into a SubscriptionProviderId
func ParseSubscriptionProviderID(input string) (*SubscriptionProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubscriptionProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubscriptionProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	return &id, nil
}

// ParseSubscriptionProviderIDInsensitively parses 'input' case-insensitively into a SubscriptionProviderId
// note: this method should only be used for API response data and not user input
func ParseSubscriptionProviderIDInsensitively(input string) (*SubscriptionProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubscriptionProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubscriptionProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	return &id, nil
}

// ValidateSubscriptionProviderID checks that 'input' can be parsed as a Subscription Provider ID
func ValidateSubscriptionProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubscriptionProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subscription Provider ID
func (id SubscriptionProviderId) ID() string {
	fmtString := "/subscriptions/%s/providers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Subscription Provider ID
func (id SubscriptionProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
	}
}

// String returns a human-readable description of this Subscription Provider ID
func (id SubscriptionProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
	}
	return fmt.Sprintf("Subscription Provider (%s)", strings.Join(components, "\n"))
}
