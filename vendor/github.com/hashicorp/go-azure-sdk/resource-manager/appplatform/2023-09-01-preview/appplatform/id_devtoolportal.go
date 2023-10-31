package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DevToolPortalId{}

// DevToolPortalId is a struct representing the Resource ID for a Dev Tool Portal
type DevToolPortalId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	DevToolPortalName string
}

// NewDevToolPortalID returns a new DevToolPortalId struct
func NewDevToolPortalID(subscriptionId string, resourceGroupName string, springName string, devToolPortalName string) DevToolPortalId {
	return DevToolPortalId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		DevToolPortalName: devToolPortalName,
	}
}

// ParseDevToolPortalID parses 'input' into a DevToolPortalId
func ParseDevToolPortalID(input string) (*DevToolPortalId, error) {
	parser := resourceids.NewParserFromResourceIdType(DevToolPortalId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DevToolPortalId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.DevToolPortalName, ok = parsed.Parsed["devToolPortalName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devToolPortalName", *parsed)
	}

	return &id, nil
}

// ParseDevToolPortalIDInsensitively parses 'input' case-insensitively into a DevToolPortalId
// note: this method should only be used for API response data and not user input
func ParseDevToolPortalIDInsensitively(input string) (*DevToolPortalId, error) {
	parser := resourceids.NewParserFromResourceIdType(DevToolPortalId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DevToolPortalId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.DevToolPortalName, ok = parsed.Parsed["devToolPortalName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "devToolPortalName", *parsed)
	}

	return &id, nil
}

// ValidateDevToolPortalID checks that 'input' can be parsed as a Dev Tool Portal ID
func ValidateDevToolPortalID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevToolPortalID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Tool Portal ID
func (id DevToolPortalId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/devToolPortals/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.DevToolPortalName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Tool Portal ID
func (id DevToolPortalId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticDevToolPortals", "devToolPortals", "devToolPortals"),
		resourceids.UserSpecifiedSegment("devToolPortalName", "devToolPortalValue"),
	}
}

// String returns a human-readable description of this Dev Tool Portal ID
func (id DevToolPortalId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Dev Tool Portal Name: %q", id.DevToolPortalName),
	}
	return fmt.Sprintf("Dev Tool Portal (%s)", strings.Join(components, "\n"))
}
