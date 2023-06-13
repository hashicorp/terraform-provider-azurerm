package managementlocks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderLockId{}

// ProviderLockId is a struct representing the Resource ID for a Provider Lock
type ProviderLockId struct {
	SubscriptionId    string
	ResourceGroupName string
	LockName          string
}

// NewProviderLockID returns a new ProviderLockId struct
func NewProviderLockID(subscriptionId string, resourceGroupName string, lockName string) ProviderLockId {
	return ProviderLockId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LockName:          lockName,
	}
}

// ParseProviderLockID parses 'input' into a ProviderLockId
func ParseProviderLockID(input string) (*ProviderLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLockId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLockId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "lockName", *parsed)
	}

	return &id, nil
}

// ParseProviderLockIDInsensitively parses 'input' case-insensitively into a ProviderLockId
// note: this method should only be used for API response data and not user input
func ParseProviderLockIDInsensitively(input string) (*ProviderLockId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLockId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLockId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LockName, ok = parsed.Parsed["lockName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "lockName", *parsed)
	}

	return &id, nil
}

// ValidateProviderLockID checks that 'input' can be parsed as a Provider Lock ID
func ValidateProviderLockID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderLockID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Lock ID
func (id ProviderLockId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Authorization/locks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LockName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Lock ID
func (id ProviderLockId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticLocks", "locks", "locks"),
		resourceids.UserSpecifiedSegment("lockName", "lockValue"),
	}
}

// String returns a human-readable description of this Provider Lock ID
func (id ProviderLockId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lock Name: %q", id.LockName),
	}
	return fmt.Sprintf("Provider Lock (%s)", strings.Join(components, "\n"))
}
