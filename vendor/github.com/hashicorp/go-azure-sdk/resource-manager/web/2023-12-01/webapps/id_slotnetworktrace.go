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
	recaser.RegisterResourceId(&SlotNetworkTraceId{})
}

var _ resourceids.ResourceId = &SlotNetworkTraceId{}

// SlotNetworkTraceId is a struct representing the Resource ID for a Slot Network Trace
type SlotNetworkTraceId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	OperationId       string
}

// NewSlotNetworkTraceID returns a new SlotNetworkTraceId struct
func NewSlotNetworkTraceID(subscriptionId string, resourceGroupName string, siteName string, slotName string, operationId string) SlotNetworkTraceId {
	return SlotNetworkTraceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		OperationId:       operationId,
	}
}

// ParseSlotNetworkTraceID parses 'input' into a SlotNetworkTraceId
func ParseSlotNetworkTraceID(input string) (*SlotNetworkTraceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotNetworkTraceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotNetworkTraceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotNetworkTraceIDInsensitively parses 'input' case-insensitively into a SlotNetworkTraceId
// note: this method should only be used for API response data and not user input
func ParseSlotNetworkTraceIDInsensitively(input string) (*SlotNetworkTraceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotNetworkTraceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotNetworkTraceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotNetworkTraceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.OperationId, ok = input.Parsed["operationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationId", input)
	}

	return nil
}

// ValidateSlotNetworkTraceID checks that 'input' can be parsed as a Slot Network Trace ID
func ValidateSlotNetworkTraceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotNetworkTraceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Network Trace ID
func (id SlotNetworkTraceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/networkTrace/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Network Trace ID
func (id SlotNetworkTraceId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticNetworkTrace", "networkTrace", "networkTrace"),
		resourceids.UserSpecifiedSegment("operationId", "operationId"),
	}
}

// String returns a human-readable description of this Slot Network Trace ID
func (id SlotNetworkTraceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Slot Network Trace (%s)", strings.Join(components, "\n"))
}
