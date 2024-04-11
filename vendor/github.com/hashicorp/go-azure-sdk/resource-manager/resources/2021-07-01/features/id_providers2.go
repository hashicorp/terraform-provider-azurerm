package features

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &Providers2Id{}

// Providers2Id is a struct representing the Resource ID for a Providers 2
type Providers2Id struct {
	SubscriptionId string
	ProviderName   string
}

// NewProviders2ID returns a new Providers2Id struct
func NewProviders2ID(subscriptionId string, providerName string) Providers2Id {
	return Providers2Id{
		SubscriptionId: subscriptionId,
		ProviderName:   providerName,
	}
}

// ParseProviders2ID parses 'input' into a Providers2Id
func ParseProviders2ID(input string) (*Providers2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2Id{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2Id{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviders2IDInsensitively parses 'input' case-insensitively into a Providers2Id
// note: this method should only be used for API response data and not user input
func ParseProviders2IDInsensitively(input string) (*Providers2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2Id{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2Id{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Providers2Id) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ProviderName, ok = input.Parsed["providerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "providerName", input)
	}

	return nil
}

// ValidateProviders2ID checks that 'input' can be parsed as a Providers 2 ID
func ValidateProviders2ID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2ID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 ID
func (id Providers2Id) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Features/providers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 ID
func (id Providers2Id) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftFeatures", "Microsoft.Features", "Microsoft.Features"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
	}
}

// String returns a human-readable description of this Providers 2 ID
func (id Providers2Id) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
	}
	return fmt.Sprintf("Providers 2 (%s)", strings.Join(components, "\n"))
}
