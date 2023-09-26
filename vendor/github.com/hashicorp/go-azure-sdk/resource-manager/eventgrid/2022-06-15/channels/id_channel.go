package channels

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ChannelId{}

// ChannelId is a struct representing the Resource ID for a Channel
type ChannelId struct {
	SubscriptionId       string
	ResourceGroupName    string
	PartnerNamespaceName string
	ChannelName          string
}

// NewChannelID returns a new ChannelId struct
func NewChannelID(subscriptionId string, resourceGroupName string, partnerNamespaceName string, channelName string) ChannelId {
	return ChannelId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		PartnerNamespaceName: partnerNamespaceName,
		ChannelName:          channelName,
	}
}

// ParseChannelID parses 'input' into a ChannelId
func ParseChannelID(input string) (*ChannelId, error) {
	parser := resourceids.NewParserFromResourceIdType(ChannelId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ChannelId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PartnerNamespaceName, ok = parsed.Parsed["partnerNamespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerNamespaceName", *parsed)
	}

	if id.ChannelName, ok = parsed.Parsed["channelName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "channelName", *parsed)
	}

	return &id, nil
}

// ParseChannelIDInsensitively parses 'input' case-insensitively into a ChannelId
// note: this method should only be used for API response data and not user input
func ParseChannelIDInsensitively(input string) (*ChannelId, error) {
	parser := resourceids.NewParserFromResourceIdType(ChannelId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ChannelId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PartnerNamespaceName, ok = parsed.Parsed["partnerNamespaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partnerNamespaceName", *parsed)
	}

	if id.ChannelName, ok = parsed.Parsed["channelName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "channelName", *parsed)
	}

	return &id, nil
}

// ValidateChannelID checks that 'input' can be parsed as a Channel ID
func ValidateChannelID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseChannelID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Channel ID
func (id ChannelId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/partnerNamespaces/%s/channels/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PartnerNamespaceName, id.ChannelName)
}

// Segments returns a slice of Resource ID Segments which comprise this Channel ID
func (id ChannelId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticPartnerNamespaces", "partnerNamespaces", "partnerNamespaces"),
		resourceids.UserSpecifiedSegment("partnerNamespaceName", "partnerNamespaceValue"),
		resourceids.StaticSegment("staticChannels", "channels", "channels"),
		resourceids.UserSpecifiedSegment("channelName", "channelValue"),
	}
}

// String returns a human-readable description of this Channel ID
func (id ChannelId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Partner Namespace Name: %q", id.PartnerNamespaceName),
		fmt.Sprintf("Channel Name: %q", id.ChannelName),
	}
	return fmt.Sprintf("Channel (%s)", strings.Join(components, "\n"))
}
