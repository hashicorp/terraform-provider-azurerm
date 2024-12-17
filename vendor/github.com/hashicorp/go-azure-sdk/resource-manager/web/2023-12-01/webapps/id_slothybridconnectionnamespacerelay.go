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
	recaser.RegisterResourceId(&SlotHybridConnectionNamespaceRelayId{})
}

var _ resourceids.ResourceId = &SlotHybridConnectionNamespaceRelayId{}

// SlotHybridConnectionNamespaceRelayId is a struct representing the Resource ID for a Slot Hybrid Connection Namespace Relay
type SlotHybridConnectionNamespaceRelayId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SiteName                      string
	SlotName                      string
	HybridConnectionNamespaceName string
	RelayName                     string
}

// NewSlotHybridConnectionNamespaceRelayID returns a new SlotHybridConnectionNamespaceRelayId struct
func NewSlotHybridConnectionNamespaceRelayID(subscriptionId string, resourceGroupName string, siteName string, slotName string, hybridConnectionNamespaceName string, relayName string) SlotHybridConnectionNamespaceRelayId {
	return SlotHybridConnectionNamespaceRelayId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SiteName:                      siteName,
		SlotName:                      slotName,
		HybridConnectionNamespaceName: hybridConnectionNamespaceName,
		RelayName:                     relayName,
	}
}

// ParseSlotHybridConnectionNamespaceRelayID parses 'input' into a SlotHybridConnectionNamespaceRelayId
func ParseSlotHybridConnectionNamespaceRelayID(input string) (*SlotHybridConnectionNamespaceRelayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotHybridConnectionNamespaceRelayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotHybridConnectionNamespaceRelayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotHybridConnectionNamespaceRelayIDInsensitively parses 'input' case-insensitively into a SlotHybridConnectionNamespaceRelayId
// note: this method should only be used for API response data and not user input
func ParseSlotHybridConnectionNamespaceRelayIDInsensitively(input string) (*SlotHybridConnectionNamespaceRelayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotHybridConnectionNamespaceRelayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotHybridConnectionNamespaceRelayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotHybridConnectionNamespaceRelayId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.HybridConnectionNamespaceName, ok = input.Parsed["hybridConnectionNamespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridConnectionNamespaceName", input)
	}

	if id.RelayName, ok = input.Parsed["relayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "relayName", input)
	}

	return nil
}

// ValidateSlotHybridConnectionNamespaceRelayID checks that 'input' can be parsed as a Slot Hybrid Connection Namespace Relay ID
func ValidateSlotHybridConnectionNamespaceRelayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotHybridConnectionNamespaceRelayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Hybrid Connection Namespace Relay ID
func (id SlotHybridConnectionNamespaceRelayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/hybridConnectionNamespaces/%s/relays/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.HybridConnectionNamespaceName, id.RelayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Hybrid Connection Namespace Relay ID
func (id SlotHybridConnectionNamespaceRelayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotName"),
		resourceids.StaticSegment("staticHybridConnectionNamespaces", "hybridConnectionNamespaces", "hybridConnectionNamespaces"),
		resourceids.UserSpecifiedSegment("hybridConnectionNamespaceName", "hybridConnectionNamespaceName"),
		resourceids.StaticSegment("staticRelays", "relays", "relays"),
		resourceids.UserSpecifiedSegment("relayName", "relayName"),
	}
}

// String returns a human-readable description of this Slot Hybrid Connection Namespace Relay ID
func (id SlotHybridConnectionNamespaceRelayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Hybrid Connection Namespace Name: %q", id.HybridConnectionNamespaceName),
		fmt.Sprintf("Relay Name: %q", id.RelayName),
	}
	return fmt.Sprintf("Slot Hybrid Connection Namespace Relay (%s)", strings.Join(components, "\n"))
}
