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
	recaser.RegisterResourceId(&SlotHostNameBindingId{})
}

var _ resourceids.ResourceId = &SlotHostNameBindingId{}

// SlotHostNameBindingId is a struct representing the Resource ID for a Slot Host Name Binding
type SlotHostNameBindingId struct {
	SubscriptionId      string
	ResourceGroupName   string
	SiteName            string
	SlotName            string
	HostNameBindingName string
}

// NewSlotHostNameBindingID returns a new SlotHostNameBindingId struct
func NewSlotHostNameBindingID(subscriptionId string, resourceGroupName string, siteName string, slotName string, hostNameBindingName string) SlotHostNameBindingId {
	return SlotHostNameBindingId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		SiteName:            siteName,
		SlotName:            slotName,
		HostNameBindingName: hostNameBindingName,
	}
}

// ParseSlotHostNameBindingID parses 'input' into a SlotHostNameBindingId
func ParseSlotHostNameBindingID(input string) (*SlotHostNameBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotHostNameBindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotHostNameBindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotHostNameBindingIDInsensitively parses 'input' case-insensitively into a SlotHostNameBindingId
// note: this method should only be used for API response data and not user input
func ParseSlotHostNameBindingIDInsensitively(input string) (*SlotHostNameBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotHostNameBindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotHostNameBindingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotHostNameBindingId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.HostNameBindingName, ok = input.Parsed["hostNameBindingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostNameBindingName", input)
	}

	return nil
}

// ValidateSlotHostNameBindingID checks that 'input' can be parsed as a Slot Host Name Binding ID
func ValidateSlotHostNameBindingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotHostNameBindingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Host Name Binding ID
func (id SlotHostNameBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/hostNameBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.HostNameBindingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Host Name Binding ID
func (id SlotHostNameBindingId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticHostNameBindings", "hostNameBindings", "hostNameBindings"),
		resourceids.UserSpecifiedSegment("hostNameBindingName", "hostNameBindingName"),
	}
}

// String returns a human-readable description of this Slot Host Name Binding ID
func (id SlotHostNameBindingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Host Name Binding Name: %q", id.HostNameBindingName),
	}
	return fmt.Sprintf("Slot Host Name Binding (%s)", strings.Join(components, "\n"))
}
