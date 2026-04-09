package inboundendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InboundEndpointId{})
}

var _ resourceids.ResourceId = &InboundEndpointId{}

// InboundEndpointId is a struct representing the Resource ID for a Inbound Endpoint
type InboundEndpointId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DnsResolverName     string
	InboundEndpointName string
}

// NewInboundEndpointID returns a new InboundEndpointId struct
func NewInboundEndpointID(subscriptionId string, resourceGroupName string, dnsResolverName string, inboundEndpointName string) InboundEndpointId {
	return InboundEndpointId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DnsResolverName:     dnsResolverName,
		InboundEndpointName: inboundEndpointName,
	}
}

// ParseInboundEndpointID parses 'input' into a InboundEndpointId
func ParseInboundEndpointID(input string) (*InboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InboundEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InboundEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInboundEndpointIDInsensitively parses 'input' case-insensitively into a InboundEndpointId
// note: this method should only be used for API response data and not user input
func ParseInboundEndpointIDInsensitively(input string) (*InboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InboundEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InboundEndpointId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InboundEndpointId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.InboundEndpointName, ok = input.Parsed["inboundEndpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "inboundEndpointName", input)
	}

	return nil
}

// ValidateInboundEndpointID checks that 'input' can be parsed as a Inbound Endpoint ID
func ValidateInboundEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInboundEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Inbound Endpoint ID
func (id InboundEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsResolvers/%s/inboundEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName, id.InboundEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Inbound Endpoint ID
func (id InboundEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsResolvers", "dnsResolvers", "dnsResolvers"),
		resourceids.UserSpecifiedSegment("dnsResolverName", "dnsResolverName"),
		resourceids.StaticSegment("staticInboundEndpoints", "inboundEndpoints", "inboundEndpoints"),
		resourceids.UserSpecifiedSegment("inboundEndpointName", "inboundEndpointName"),
	}
}

// String returns a human-readable description of this Inbound Endpoint ID
func (id InboundEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Resolver Name: %q", id.DnsResolverName),
		fmt.Sprintf("Inbound Endpoint Name: %q", id.InboundEndpointName),
	}
	return fmt.Sprintf("Inbound Endpoint (%s)", strings.Join(components, "\n"))
}
