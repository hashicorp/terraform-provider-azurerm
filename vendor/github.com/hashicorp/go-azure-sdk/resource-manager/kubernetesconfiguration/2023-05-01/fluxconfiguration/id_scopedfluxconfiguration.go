package fluxconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedFluxConfigurationId{}

// ScopedFluxConfigurationId is a struct representing the Resource ID for a Scoped Flux Configuration
type ScopedFluxConfigurationId struct {
	Scope                 string
	FluxConfigurationName string
}

// NewScopedFluxConfigurationID returns a new ScopedFluxConfigurationId struct
func NewScopedFluxConfigurationID(scope string, fluxConfigurationName string) ScopedFluxConfigurationId {
	return ScopedFluxConfigurationId{
		Scope:                 scope,
		FluxConfigurationName: fluxConfigurationName,
	}
}

// ParseScopedFluxConfigurationID parses 'input' into a ScopedFluxConfigurationId
func ParseScopedFluxConfigurationID(input string) (*ScopedFluxConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedFluxConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedFluxConfigurationId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.FluxConfigurationName, ok = parsed.Parsed["fluxConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluxConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseScopedFluxConfigurationIDInsensitively parses 'input' case-insensitively into a ScopedFluxConfigurationId
// note: this method should only be used for API response data and not user input
func ParseScopedFluxConfigurationIDInsensitively(input string) (*ScopedFluxConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedFluxConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedFluxConfigurationId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.FluxConfigurationName, ok = parsed.Parsed["fluxConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluxConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateScopedFluxConfigurationID checks that 'input' can be parsed as a Scoped Flux Configuration ID
func ValidateScopedFluxConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedFluxConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Flux Configuration ID
func (id ScopedFluxConfigurationId) ID() string {
	fmtString := "/%s/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.FluxConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Flux Configuration ID
func (id ScopedFluxConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKubernetesConfiguration", "Microsoft.KubernetesConfiguration", "Microsoft.KubernetesConfiguration"),
		resourceids.StaticSegment("staticFluxConfigurations", "fluxConfigurations", "fluxConfigurations"),
		resourceids.UserSpecifiedSegment("fluxConfigurationName", "fluxConfigurationValue"),
	}
}

// String returns a human-readable description of this Scoped Flux Configuration ID
func (id ScopedFluxConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Flux Configuration Name: %q", id.FluxConfigurationName),
	}
	return fmt.Sprintf("Scoped Flux Configuration (%s)", strings.Join(components, "\n"))
}
