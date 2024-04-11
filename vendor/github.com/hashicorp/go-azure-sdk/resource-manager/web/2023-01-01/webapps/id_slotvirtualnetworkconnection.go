package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SlotVirtualNetworkConnectionId{}

// SlotVirtualNetworkConnectionId is a struct representing the Resource ID for a Slot Virtual Network Connection
type SlotVirtualNetworkConnectionId struct {
	SubscriptionId               string
	ResourceGroupName            string
	SiteName                     string
	SlotName                     string
	VirtualNetworkConnectionName string
}

// NewSlotVirtualNetworkConnectionID returns a new SlotVirtualNetworkConnectionId struct
func NewSlotVirtualNetworkConnectionID(subscriptionId string, resourceGroupName string, siteName string, slotName string, virtualNetworkConnectionName string) SlotVirtualNetworkConnectionId {
	return SlotVirtualNetworkConnectionId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		SiteName:                     siteName,
		SlotName:                     slotName,
		VirtualNetworkConnectionName: virtualNetworkConnectionName,
	}
}

// ParseSlotVirtualNetworkConnectionID parses 'input' into a SlotVirtualNetworkConnectionId
func ParseSlotVirtualNetworkConnectionID(input string) (*SlotVirtualNetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotVirtualNetworkConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotVirtualNetworkConnectionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotVirtualNetworkConnectionIDInsensitively parses 'input' case-insensitively into a SlotVirtualNetworkConnectionId
// note: this method should only be used for API response data and not user input
func ParseSlotVirtualNetworkConnectionIDInsensitively(input string) (*SlotVirtualNetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotVirtualNetworkConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotVirtualNetworkConnectionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotVirtualNetworkConnectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.VirtualNetworkConnectionName, ok = input.Parsed["virtualNetworkConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkConnectionName", input)
	}

	return nil
}

// ValidateSlotVirtualNetworkConnectionID checks that 'input' can be parsed as a Slot Virtual Network Connection ID
func ValidateSlotVirtualNetworkConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotVirtualNetworkConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Virtual Network Connection ID
func (id SlotVirtualNetworkConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/virtualNetworkConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.VirtualNetworkConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Virtual Network Connection ID
func (id SlotVirtualNetworkConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteValue"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotValue"),
		resourceids.StaticSegment("staticVirtualNetworkConnections", "virtualNetworkConnections", "virtualNetworkConnections"),
		resourceids.UserSpecifiedSegment("virtualNetworkConnectionName", "virtualNetworkConnectionValue"),
	}
}

// String returns a human-readable description of this Slot Virtual Network Connection ID
func (id SlotVirtualNetworkConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Virtual Network Connection Name: %q", id.VirtualNetworkConnectionName),
	}
	return fmt.Sprintf("Slot Virtual Network Connection (%s)", strings.Join(components, "\n"))
}
