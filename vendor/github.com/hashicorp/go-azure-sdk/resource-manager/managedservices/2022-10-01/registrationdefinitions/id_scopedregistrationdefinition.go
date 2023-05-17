package registrationdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedRegistrationDefinitionId{}

// ScopedRegistrationDefinitionId is a struct representing the Resource ID for a Scoped Registration Definition
type ScopedRegistrationDefinitionId struct {
	Scope                    string
	RegistrationDefinitionId string
}

// NewScopedRegistrationDefinitionID returns a new ScopedRegistrationDefinitionId struct
func NewScopedRegistrationDefinitionID(scope string, registrationDefinitionId string) ScopedRegistrationDefinitionId {
	return ScopedRegistrationDefinitionId{
		Scope:                    scope,
		RegistrationDefinitionId: registrationDefinitionId,
	}
}

// ParseScopedRegistrationDefinitionID parses 'input' into a ScopedRegistrationDefinitionId
func ParseScopedRegistrationDefinitionID(input string) (*ScopedRegistrationDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRegistrationDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRegistrationDefinitionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RegistrationDefinitionId, ok = parsed.Parsed["registrationDefinitionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registrationDefinitionId", *parsed)
	}

	return &id, nil
}

// ParseScopedRegistrationDefinitionIDInsensitively parses 'input' case-insensitively into a ScopedRegistrationDefinitionId
// note: this method should only be used for API response data and not user input
func ParseScopedRegistrationDefinitionIDInsensitively(input string) (*ScopedRegistrationDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedRegistrationDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedRegistrationDefinitionId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.RegistrationDefinitionId, ok = parsed.Parsed["registrationDefinitionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registrationDefinitionId", *parsed)
	}

	return &id, nil
}

// ValidateScopedRegistrationDefinitionID checks that 'input' can be parsed as a Scoped Registration Definition ID
func ValidateScopedRegistrationDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedRegistrationDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Registration Definition ID
func (id ScopedRegistrationDefinitionId) ID() string {
	fmtString := "/%s/providers/Microsoft.ManagedServices/registrationDefinitions/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.RegistrationDefinitionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Registration Definition ID
func (id ScopedRegistrationDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagedServices", "Microsoft.ManagedServices", "Microsoft.ManagedServices"),
		resourceids.StaticSegment("staticRegistrationDefinitions", "registrationDefinitions", "registrationDefinitions"),
		resourceids.UserSpecifiedSegment("registrationDefinitionId", "registrationDefinitionIdValue"),
	}
}

// String returns a human-readable description of this Scoped Registration Definition ID
func (id ScopedRegistrationDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Registration Definition: %q", id.RegistrationDefinitionId),
	}
	return fmt.Sprintf("Scoped Registration Definition (%s)", strings.Join(components, "\n"))
}
