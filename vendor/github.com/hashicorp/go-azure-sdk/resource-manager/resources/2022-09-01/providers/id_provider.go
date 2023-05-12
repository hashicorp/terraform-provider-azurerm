package providers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderId{}

// ProviderId is a struct representing the Resource ID for a Provider
type ProviderId struct {
	ProviderName string
}

// NewProviderID returns a new ProviderId struct
func NewProviderID(providerName string) ProviderId {
	return ProviderId{
		ProviderName: providerName,
	}
}

// ParseProviderID parses 'input' into a ProviderId
func ParseProviderID(input string) (*ProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderId{}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	return &id, nil
}

// ParseProviderIDInsensitively parses 'input' case-insensitively into a ProviderId
// note: this method should only be used for API response data and not user input
func ParseProviderIDInsensitively(input string) (*ProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderId{}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	return &id, nil
}

// ValidateProviderID checks that 'input' can be parsed as a Provider ID
func ValidateProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider ID
func (id ProviderId) ID() string {
	fmtString := "/providers/%s"
	return fmt.Sprintf(fmtString, id.ProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider ID
func (id ProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
	}
}

// String returns a human-readable description of this Provider ID
func (id ProviderId) String() string {
	components := []string{
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
	}
	return fmt.Sprintf("Provider (%s)", strings.Join(components, "\n"))
}
