package expressroutelinks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LinkId{}

// LinkId is a struct representing the Resource ID for a Link
type LinkId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ExpressRoutePortName string
	LinkName             string
}

// NewLinkID returns a new LinkId struct
func NewLinkID(subscriptionId string, resourceGroupName string, expressRoutePortName string, linkName string) LinkId {
	return LinkId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ExpressRoutePortName: expressRoutePortName,
		LinkName:             linkName,
	}
}

// ParseLinkID parses 'input' into a LinkId
func ParseLinkID(input string) (*LinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(LinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LinkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRoutePortName, ok = parsed.Parsed["expressRoutePortName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortName", *parsed)
	}

	if id.LinkName, ok = parsed.Parsed["linkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "linkName", *parsed)
	}

	return &id, nil
}

// ParseLinkIDInsensitively parses 'input' case-insensitively into a LinkId
// note: this method should only be used for API response data and not user input
func ParseLinkIDInsensitively(input string) (*LinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(LinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LinkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRoutePortName, ok = parsed.Parsed["expressRoutePortName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortName", *parsed)
	}

	if id.LinkName, ok = parsed.Parsed["linkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "linkName", *parsed)
	}

	return &id, nil
}

// ValidateLinkID checks that 'input' can be parsed as a Link ID
func ValidateLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Link ID
func (id LinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRoutePorts/%s/links/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName, id.LinkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Link ID
func (id LinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRoutePorts", "expressRoutePorts", "expressRoutePorts"),
		resourceids.UserSpecifiedSegment("expressRoutePortName", "expressRoutePortValue"),
		resourceids.StaticSegment("staticLinks", "links", "links"),
		resourceids.UserSpecifiedSegment("linkName", "linkValue"),
	}
}

// String returns a human-readable description of this Link ID
func (id LinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Port Name: %q", id.ExpressRoutePortName),
		fmt.Sprintf("Link Name: %q", id.LinkName),
	}
	return fmt.Sprintf("Link (%s)", strings.Join(components, "\n"))
}
