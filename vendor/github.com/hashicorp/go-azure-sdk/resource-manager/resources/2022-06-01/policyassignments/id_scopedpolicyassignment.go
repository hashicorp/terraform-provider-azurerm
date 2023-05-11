package policyassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedPolicyAssignmentId{}

// ScopedPolicyAssignmentId is a struct representing the Resource ID for a Scoped Policy Assignment
type ScopedPolicyAssignmentId struct {
	Scope                string
	PolicyAssignmentName string
}

// NewScopedPolicyAssignmentID returns a new ScopedPolicyAssignmentId struct
func NewScopedPolicyAssignmentID(scope string, policyAssignmentName string) ScopedPolicyAssignmentId {
	return ScopedPolicyAssignmentId{
		Scope:                scope,
		PolicyAssignmentName: policyAssignmentName,
	}
}

// ParseScopedPolicyAssignmentID parses 'input' into a ScopedPolicyAssignmentId
func ParseScopedPolicyAssignmentID(input string) (*ScopedPolicyAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedPolicyAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedPolicyAssignmentId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.PolicyAssignmentName, ok = parsed.Parsed["policyAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policyAssignmentName", *parsed)
	}

	return &id, nil
}

// ParseScopedPolicyAssignmentIDInsensitively parses 'input' case-insensitively into a ScopedPolicyAssignmentId
// note: this method should only be used for API response data and not user input
func ParseScopedPolicyAssignmentIDInsensitively(input string) (*ScopedPolicyAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedPolicyAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedPolicyAssignmentId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.PolicyAssignmentName, ok = parsed.Parsed["policyAssignmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "policyAssignmentName", *parsed)
	}

	return &id, nil
}

// ValidateScopedPolicyAssignmentID checks that 'input' can be parsed as a Scoped Policy Assignment ID
func ValidateScopedPolicyAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedPolicyAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Policy Assignment ID
func (id ScopedPolicyAssignmentId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/policyAssignments/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.PolicyAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Policy Assignment ID
func (id ScopedPolicyAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticPolicyAssignments", "policyAssignments", "policyAssignments"),
		resourceids.UserSpecifiedSegment("policyAssignmentName", "policyAssignmentValue"),
	}
}

// String returns a human-readable description of this Scoped Policy Assignment ID
func (id ScopedPolicyAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Policy Assignment Name: %q", id.PolicyAssignmentName),
	}
	return fmt.Sprintf("Scoped Policy Assignment (%s)", strings.Join(components, "\n"))
}
