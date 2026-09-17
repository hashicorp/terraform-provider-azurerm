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
	recaser.RegisterResourceId(&SlotHybridConnectionId{})
}

var _ resourceids.ResourceId = &SlotHybridConnectionId{}

// SlotHybridConnectionId is a struct representing the Resource ID for a Slot Hybrid Connection
type SlotHybridConnectionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	SiteName             string
	SlotName             string
	HybridConnectionName string
}

// NewSlotHybridConnectionID returns a new SlotHybridConnectionId struct
func NewSlotHybridConnectionID(subscriptionId string, resourceGroupName string, siteName string, slotName string, hybridConnectionName string) SlotHybridConnectionId {
	return SlotHybridConnectionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		SiteName:             siteName,
		SlotName:             slotName,
		HybridConnectionName: hybridConnectionName,
	}
}

// ParseSlotHybridConnectionID parses 'input' into a SlotHybridConnectionId
func ParseSlotHybridConnectionID(input string) (*SlotHybridConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotHybridConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotHybridConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotHybridConnectionIDInsensitively parses 'input' case-insensitively into a SlotHybridConnectionId
// note: this method should only be used for API response data and not user input
func ParseSlotHybridConnectionIDInsensitively(input string) (*SlotHybridConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotHybridConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotHybridConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotHybridConnectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.HybridConnectionName, ok = input.Parsed["hybridConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hybridConnectionName", input)
	}

	return nil
}

// ValidateSlotHybridConnectionID checks that 'input' can be parsed as a Slot Hybrid Connection ID
func ValidateSlotHybridConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotHybridConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Hybrid Connection ID
func (id SlotHybridConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/hybridConnection/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.HybridConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Hybrid Connection ID
func (id SlotHybridConnectionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticHybridConnection", "hybridConnection", "hybridConnection"),
		resourceids.UserSpecifiedSegment("hybridConnectionName", "hybridConnectionName"),
	}
}

// String returns a human-readable description of this Slot Hybrid Connection ID
func (id SlotHybridConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Hybrid Connection Name: %q", id.HybridConnectionName),
	}
	return fmt.Sprintf("Slot Hybrid Connection (%s)", strings.Join(components, "\n"))
}
