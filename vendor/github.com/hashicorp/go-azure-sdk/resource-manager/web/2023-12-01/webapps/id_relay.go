package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RelayId{})
}

var _ resourceids.ResourceId = &RelayId{}

// RelayId is a struct representing the Resource ID for a Relay
type RelayId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SiteName                      string
	HybridConnectionNamespaceName string
	RelayName                     string
}

// NewRelayID returns a new RelayId struct
func NewRelayID(subscriptionId string, resourceGroupName string, siteName string, hybridConnectionNamespaceName string, relayName string) RelayId {
	return RelayId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SiteName:                      siteName,
		HybridConnectionNamespaceName: hybridConnectionNamespaceName,
		RelayName:                     relayName,
	}
}

// ParseRelayID parses 'input' into a RelayId
func ParseRelayID(input string) (*RelayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RelayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RelayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRelayIDInsensitively parses 'input' case-insensitively into a RelayId
// note: this method should only be used for API response data and not user input
func ParseRelayIDInsensitively(input string) (*RelayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RelayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RelayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RelayId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.HybridConnectionNamespaceName, ok = input.Parsed["hybridConnectionNamespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridConnectionNamespaceName", input)
	}

	if id.RelayName, ok = input.Parsed["relayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "relayName", input)
	}

	return nil
}

// ValidateRelayID checks that 'input' can be parsed as a Relay ID
func ValidateRelayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRelayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Relay ID
func (id RelayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hybridConnectionNamespaces/%s/relays/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Relay ID
func (id RelayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticHybridConnectionNamespaces", "hybridConnectionNamespaces", "hybridConnectionNamespaces"),
		resourceids.UserSpecifiedSegment("hybridConnectionNamespaceName", "hybridConnectionNamespaceName"),
		resourceids.StaticSegment("staticRelays", "relays", "relays"),
		resourceids.UserSpecifiedSegment("relayName", "relayName"),
	}
}

// String returns a human-readable description of this Relay ID
func (id RelayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Hybrid Connection Namespace Name: %q", id.HybridConnectionNamespaceName),
		fmt.Sprintf("Relay Name: %q", id.RelayName),
	}
	return fmt.Sprintf("Relay (%s)", strings.Join(components, "\n"))
}
