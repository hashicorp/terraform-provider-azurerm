package managedhsms

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ManagedHSMId{}

// ManagedHSMId is a struct representing the Resource ID for a Managed H S M
type ManagedHSMId struct {
	SubscriptionId    string
	ResourceGroupName string
	ManagedHSMName    string
}

// NewManagedHSMID returns a new ManagedHSMId struct
func NewManagedHSMID(subscriptionId string, resourceGroupName string, managedHSMName string) ManagedHSMId {
	return ManagedHSMId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ManagedHSMName:    managedHSMName,
	}
}

// ParseManagedHSMID parses 'input' into a ManagedHSMId
func ParseManagedHSMID(input string) (*ManagedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedHSMId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedHSMId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedHSMName, ok = parsed.Parsed["managedHSMName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedHSMName", *parsed)
	}

	return &id, nil
}

// ParseManagedHSMIDInsensitively parses 'input' case-insensitively into a ManagedHSMId
// note: this method should only be used for API response data and not user input
func ParseManagedHSMIDInsensitively(input string) (*ManagedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedHSMId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedHSMId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedHSMName, ok = parsed.Parsed["managedHSMName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedHSMName", *parsed)
	}

	return &id, nil
}

// ValidateManagedHSMID checks that 'input' can be parsed as a Managed H S M ID
func ValidateManagedHSMID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedHSMID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed H S M ID
func (id ManagedHSMId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedHSMName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed H S M ID
func (id ManagedHSMId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticManagedHSMs", "managedHSMs", "managedHSMs"),
		resourceids.UserSpecifiedSegment("managedHSMName", "managedHSMValue"),
	}
}

// String returns a human-readable description of this Managed H S M ID
func (id ManagedHSMId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed H S M Name: %q", id.ManagedHSMName),
	}
	return fmt.Sprintf("Managed H S M (%s)", strings.Join(components, "\n"))
}
