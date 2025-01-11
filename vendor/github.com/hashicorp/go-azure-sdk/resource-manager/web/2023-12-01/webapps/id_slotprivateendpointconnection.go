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
	recaser.RegisterResourceId(&SlotPrivateEndpointConnectionId{})
}

var _ resourceids.ResourceId = &SlotPrivateEndpointConnectionId{}

// SlotPrivateEndpointConnectionId is a struct representing the Resource ID for a Slot Private Endpoint Connection
type SlotPrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroupName             string
	SiteName                      string
	SlotName                      string
	PrivateEndpointConnectionName string
}

// NewSlotPrivateEndpointConnectionID returns a new SlotPrivateEndpointConnectionId struct
func NewSlotPrivateEndpointConnectionID(subscriptionId string, resourceGroupName string, siteName string, slotName string, privateEndpointConnectionName string) SlotPrivateEndpointConnectionId {
	return SlotPrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		SiteName:                      siteName,
		SlotName:                      slotName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

// ParseSlotPrivateEndpointConnectionID parses 'input' into a SlotPrivateEndpointConnectionId
func ParseSlotPrivateEndpointConnectionID(input string) (*SlotPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotPrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a SlotPrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParseSlotPrivateEndpointConnectionIDInsensitively(input string) (*SlotPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotPrivateEndpointConnectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PrivateEndpointConnectionName, ok = input.Parsed["privateEndpointConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", input)
	}

	return nil
}

// ValidateSlotPrivateEndpointConnectionID checks that 'input' can be parsed as a Slot Private Endpoint Connection ID
func ValidateSlotPrivateEndpointConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotPrivateEndpointConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Private Endpoint Connection ID
func (id SlotPrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.PrivateEndpointConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Private Endpoint Connection ID
func (id SlotPrivateEndpointConnectionId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionName"),
	}
}

// String returns a human-readable description of this Slot Private Endpoint Connection ID
func (id SlotPrivateEndpointConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Private Endpoint Connection Name: %q", id.PrivateEndpointConnectionName),
	}
	return fmt.Sprintf("Slot Private Endpoint Connection (%s)", strings.Join(components, "\n"))
}
