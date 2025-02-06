package pricings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PricingId{})
}

var _ resourceids.ResourceId = &PricingId{}

// PricingId is a struct representing the Resource ID for a Pricing
type PricingId struct {
	SubscriptionId string
	PricingName    string
}

// NewPricingID returns a new PricingId struct
func NewPricingID(subscriptionId string, pricingName string) PricingId {
	return PricingId{
		SubscriptionId: subscriptionId,
		PricingName:    pricingName,
	}
}

// ParsePricingID parses 'input' into a PricingId
func ParsePricingID(input string) (*PricingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PricingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PricingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePricingIDInsensitively parses 'input' case-insensitively into a PricingId
// note: this method should only be used for API response data and not user input
func ParsePricingIDInsensitively(input string) (*PricingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PricingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PricingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PricingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.PricingName, ok = input.Parsed["pricingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "pricingName", input)
	}

	return nil
}

// ValidatePricingID checks that 'input' can be parsed as a Pricing ID
func ValidatePricingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePricingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Pricing ID
func (id PricingId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/pricings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PricingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Pricing ID
func (id PricingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurity", "Microsoft.Security", "Microsoft.Security"),
		resourceids.StaticSegment("staticPricings", "pricings", "pricings"),
		resourceids.UserSpecifiedSegment("pricingName", "pricingName"),
	}
}

// String returns a human-readable description of this Pricing ID
func (id PricingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Pricing Name: %q", id.PricingName),
	}
	return fmt.Sprintf("Pricing (%s)", strings.Join(components, "\n"))
}
