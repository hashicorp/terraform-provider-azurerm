package outboundendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OutboundEndpointId{})
}

var _ resourceids.ResourceId = &OutboundEndpointId{}

// OutboundEndpointId is a struct representing the Resource ID for a Outbound Endpoint
type OutboundEndpointId struct {
	SubscriptionId       string
	ResourceGroupName    string
	DnsResolverName      string
	OutboundEndpointName string
}

// NewOutboundEndpointID returns a new OutboundEndpointId struct
func NewOutboundEndpointID(subscriptionId string, resourceGroupName string, dnsResolverName string, outboundEndpointName string) OutboundEndpointId {
	return OutboundEndpointId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		DnsResolverName:      dnsResolverName,
		OutboundEndpointName: outboundEndpointName,
	}
}

// ParseOutboundEndpointID parses 'input' into a OutboundEndpointId
func ParseOutboundEndpointID(input string) (*OutboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OutboundEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OutboundEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOutboundEndpointIDInsensitively parses 'input' case-insensitively into a OutboundEndpointId
// note: this method should only be used for API response data and not user input
func ParseOutboundEndpointIDInsensitively(input string) (*OutboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OutboundEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OutboundEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OutboundEndpointId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DnsResolverName, ok = input.Parsed["dnsResolverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsResolverName", input)
	}

	if id.OutboundEndpointName, ok = input.Parsed["outboundEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "outboundEndpointName", input)
	}

	return nil
}

// ValidateOutboundEndpointID checks that 'input' can be parsed as a Outbound Endpoint ID
func ValidateOutboundEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOutboundEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Outbound Endpoint ID
func (id OutboundEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsResolvers/%s/outboundEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName, id.OutboundEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Outbound Endpoint ID
func (id OutboundEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsResolvers", "dnsResolvers", "dnsResolvers"),
		resourceids.UserSpecifiedSegment("dnsResolverName", "dnsResolverName"),
		resourceids.StaticSegment("staticOutboundEndpoints", "outboundEndpoints", "outboundEndpoints"),
		resourceids.UserSpecifiedSegment("outboundEndpointName", "outboundEndpointName"),
	}
}

// String returns a human-readable description of this Outbound Endpoint ID
func (id OutboundEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Resolver Name: %q", id.DnsResolverName),
		fmt.Sprintf("Outbound Endpoint Name: %q", id.OutboundEndpointName),
	}
	return fmt.Sprintf("Outbound Endpoint (%s)", strings.Join(components, "\n"))
}
