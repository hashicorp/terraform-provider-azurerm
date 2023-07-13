package dedicatedhsms

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DedicatedHSMId{}

// DedicatedHSMId is a struct representing the Resource ID for a Dedicated H S M
type DedicatedHSMId struct {
	SubscriptionId    string
	ResourceGroupName string
	DedicatedHSMName  string
}

// NewDedicatedHSMID returns a new DedicatedHSMId struct
func NewDedicatedHSMID(subscriptionId string, resourceGroupName string, dedicatedHSMName string) DedicatedHSMId {
	return DedicatedHSMId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DedicatedHSMName:  dedicatedHSMName,
	}
}

// ParseDedicatedHSMID parses 'input' into a DedicatedHSMId
func ParseDedicatedHSMID(input string) (*DedicatedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(DedicatedHSMId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DedicatedHSMId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DedicatedHSMName, ok = parsed.Parsed["dedicatedHSMName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dedicatedHSMName", *parsed)
	}

	return &id, nil
}

// ParseDedicatedHSMIDInsensitively parses 'input' case-insensitively into a DedicatedHSMId
// note: this method should only be used for API response data and not user input
func ParseDedicatedHSMIDInsensitively(input string) (*DedicatedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(DedicatedHSMId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DedicatedHSMId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DedicatedHSMName, ok = parsed.Parsed["dedicatedHSMName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dedicatedHSMName", *parsed)
	}

	return &id, nil
}

// ValidateDedicatedHSMID checks that 'input' can be parsed as a Dedicated H S M ID
func ValidateDedicatedHSMID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDedicatedHSMID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dedicated H S M ID
func (id DedicatedHSMId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HardwareSecurityModules/dedicatedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DedicatedHSMName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dedicated H S M ID
func (id DedicatedHSMId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHardwareSecurityModules", "Microsoft.HardwareSecurityModules", "Microsoft.HardwareSecurityModules"),
		resourceids.StaticSegment("staticDedicatedHSMs", "dedicatedHSMs", "dedicatedHSMs"),
		resourceids.UserSpecifiedSegment("dedicatedHSMName", "dedicatedHSMValue"),
	}
}

// String returns a human-readable description of this Dedicated H S M ID
func (id DedicatedHSMId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dedicated H S M Name: %q", id.DedicatedHSMName),
	}
	return fmt.Sprintf("Dedicated H S M (%s)", strings.Join(components, "\n"))
}
