package subscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderSubscriptionId{})
}

var _ resourceids.ResourceId = &ProviderSubscriptionId{}

// ProviderSubscriptionId is a struct representing the Resource ID for a Provider Subscription
type ProviderSubscriptionId struct {
	SubscriptionId string
}

// NewProviderSubscriptionID returns a new ProviderSubscriptionId struct
func NewProviderSubscriptionID(subscriptionId string) ProviderSubscriptionId {
	return ProviderSubscriptionId{
		SubscriptionId: subscriptionId,
	}
}

// ParseProviderSubscriptionID parses 'input' into a ProviderSubscriptionId
func ParseProviderSubscriptionID(input string) (*ProviderSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderSubscriptionIDInsensitively parses 'input' case-insensitively into a ProviderSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseProviderSubscriptionIDInsensitively(input string) (*ProviderSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderSubscriptionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderSubscriptionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	return nil
}

// ValidateProviderSubscriptionID checks that 'input' can be parsed as a Provider Subscription ID
func ValidateProviderSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Subscription ID
func (id ProviderSubscriptionId) ID() string {
	fmtString := "/providers/Microsoft.Subscription/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Subscription ID
func (id ProviderSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSubscription", "Microsoft.Subscription", "Microsoft.Subscription"),
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
	}
}

// String returns a human-readable description of this Provider Subscription ID
func (id ProviderSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
	}
	return fmt.Sprintf("Provider Subscription (%s)", strings.Join(components, "\n"))
}
