package managementlocks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LockId{}

// LockId is a struct representing the Resource ID for a Lock
type LockId struct {
	SubscriptionId string
	LockName       string
}

// NewLockID returns a new LockId struct
func NewLockID(subscriptionId string, lockName string) LockId {
	return LockId{
		SubscriptionId: subscriptionId,
		LockName:       lockName,
	}
}

// ParseLockID parses 'input' into a LockId
func ParseLockID(input string) (*LockId, error) {
	parser := resourceids.NewParserFromResourceIdType(LockId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LockId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "lockName", *parsed)
	}

	return &id, nil
}

// ParseLockIDInsensitively parses 'input' case-insensitively into a LockId
// note: this method should only be used for API response data and not user input
func ParseLockIDInsensitively(input string) (*LockId, error) {
	parser := resourceids.NewParserFromResourceIdType(LockId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LockId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "lockName", *parsed)
	}

	return &id, nil
}

// ValidateLockID checks that 'input' can be parsed as a Lock ID
func ValidateLockID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLockID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Lock ID
func (id LockId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Authorization/locks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LockName)
}

// Segments returns a slice of Resource ID Segments which comprise this Lock ID
func (id LockId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticLocks", "locks", "locks"),
		resourceids.UserSpecifiedSegment("lockName", "lockValue"),
	}
}

// String returns a human-readable description of this Lock ID
func (id LockId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Lock Name: %q", id.LockName),
	}
	return fmt.Sprintf("Lock (%s)", strings.Join(components, "\n"))
}
