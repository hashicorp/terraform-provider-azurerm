package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DevToolPortalId{})
}

var _ resourceids.ResourceId = &DevToolPortalId{}

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
	parser := resourceids.NewParserFromResourceIdType(&DevToolPortalId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevToolPortalId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDevToolPortalIDInsensitively parses 'input' case-insensitively into a DevToolPortalId
// note: this method should only be used for API response data and not user input
func ParseDevToolPortalIDInsensitively(input string) (*DevToolPortalId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevToolPortalId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevToolPortalId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DevToolPortalId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.DevToolPortalName, ok = input.Parsed["devToolPortalName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devToolPortalName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticDevToolPortals", "devToolPortals", "devToolPortals"),
		resourceids.UserSpecifiedSegment("devToolPortalName", "devToolPortalName"),
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
