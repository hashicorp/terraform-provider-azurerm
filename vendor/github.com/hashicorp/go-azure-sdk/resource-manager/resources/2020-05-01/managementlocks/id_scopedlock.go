package managementlocks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopedLockId{}

// ScopedLockId is a struct representing the Resource ID for a Scoped Lock
type ScopedLockId struct {
	Scope    string
	LockName string
}

// NewScopedLockID returns a new ScopedLockId struct
func NewScopedLockID(scope string, lockName string) ScopedLockId {
	return ScopedLockId{
		Scope:    scope,
		LockName: lockName,
	}
}

// ParseScopedLockID parses 'input' into a ScopedLockId
func ParseScopedLockID(input string) (*ScopedLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedLockId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedLockId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "lockName", *parsed)
	}

	return &id, nil
}

// ParseScopedLockIDInsensitively parses 'input' case-insensitively into a ScopedLockId
// note: this method should only be used for API response data and not user input
func ParseScopedLockIDInsensitively(input string) (*ScopedLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopedLockId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopedLockId{}

	if id.Scope, ok = parsed.Parsed["scope"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scope", *parsed)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "lockName", *parsed)
	}

	return &id, nil
}

// ValidateScopedLockID checks that 'input' can be parsed as a Scoped Lock ID
func ValidateScopedLockID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedLockID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Lock ID
func (id ScopedLockId) ID() string {
	fmtString := "/%s/providers/Microsoft.Authorization/locks/%s"
	return fmt.Sprintf(fmtString, strings.TrimPrefix(id.Scope, "/"), id.LockName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Lock ID
func (id ScopedLockId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.ScopeSegment("scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticLocks", "locks", "locks"),
		resourceids.UserSpecifiedSegment("lockName", "lockValue"),
	}
}

// String returns a human-readable description of this Scoped Lock ID
func (id ScopedLockId) String() string {
	components := []string{
		fmt.Sprintf("Scope: %q", id.Scope),
		fmt.Sprintf("Lock Name: %q", id.LockName),
	}
	return fmt.Sprintf("Scoped Lock (%s)", strings.Join(components, "\n"))
}
