package peerexpressroutecircuitconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PeerConnectionId{})
}

var _ resourceids.ResourceId = &PeerConnectionId{}

// PeerConnectionId is a struct representing the Resource ID for a Peer Connection
type PeerConnectionId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ExpressRouteCircuitName string
	PeeringName             string
	PeerConnectionName      string
}

// NewPeerConnectionID returns a new PeerConnectionId struct
func NewPeerConnectionID(subscriptionId string, resourceGroupName string, expressRouteCircuitName string, peeringName string, peerConnectionName string) PeerConnectionId {
	return PeerConnectionId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
		PeerConnectionName:      peerConnectionName,
	}
}

// ParsePeerConnectionID parses 'input' into a PeerConnectionId
func ParsePeerConnectionID(input string) (*PeerConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeerConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeerConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePeerConnectionIDInsensitively parses 'input' case-insensitively into a PeerConnectionId
// note: this method should only be used for API response data and not user input
func ParsePeerConnectionIDInsensitively(input string) (*PeerConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PeerConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PeerConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PeerConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteCircuitName, ok = input.Parsed["expressRouteCircuitName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteCircuitName", input)
	}

	if id.PeeringName, ok = input.Parsed["peeringName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "peeringName", input)
	}

	if id.PeerConnectionName, ok = input.Parsed["peerConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "peerConnectionName", input)
	}

	return nil
}

// ValidatePeerConnectionID checks that 'input' can be parsed as a Peer Connection ID
func ValidatePeerConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePeerConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Peer Connection ID
func (id PeerConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/peerConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName, id.PeerConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Peer Connection ID
func (id PeerConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteCircuits", "expressRouteCircuits", "expressRouteCircuits"),
		resourceids.UserSpecifiedSegment("expressRouteCircuitName", "expressRouteCircuitName"),
		resourceids.StaticSegment("staticPeerings", "peerings", "peerings"),
		resourceids.UserSpecifiedSegment("peeringName", "peeringName"),
		resourceids.StaticSegment("staticPeerConnections", "peerConnections", "peerConnections"),
		resourceids.UserSpecifiedSegment("peerConnectionName", "peerConnectionName"),
	}
}

// String returns a human-readable description of this Peer Connection ID
func (id PeerConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Circuit Name: %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Peering Name: %q", id.PeeringName),
		fmt.Sprintf("Peer Connection Name: %q", id.PeerConnectionName),
	}
	return fmt.Sprintf("Peer Connection (%s)", strings.Join(components, "\n"))
}
