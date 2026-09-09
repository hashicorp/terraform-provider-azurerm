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
	recaser.RegisterResourceId(&SlotWebJobId{})
}

var _ resourceids.ResourceId = &SlotWebJobId{}

// SlotWebJobId is a struct representing the Resource ID for a Slot Web Job
type SlotWebJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	WebJobName        string
}

// NewSlotWebJobID returns a new SlotWebJobId struct
func NewSlotWebJobID(subscriptionId string, resourceGroupName string, siteName string, slotName string, webJobName string) SlotWebJobId {
	return SlotWebJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		WebJobName:        webJobName,
	}
}

// ParseSlotWebJobID parses 'input' into a SlotWebJobId
func ParseSlotWebJobID(input string) (*SlotWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotWebJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotWebJobIDInsensitively parses 'input' case-insensitively into a SlotWebJobId
// note: this method should only be used for API response data and not user input
func ParseSlotWebJobIDInsensitively(input string) (*SlotWebJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotWebJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotWebJobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotWebJobId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.WebJobName, ok = input.Parsed["webJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "webJobName", input)
	}

	return nil
}

// ValidateSlotWebJobID checks that 'input' can be parsed as a Slot Web Job ID
func ValidateSlotWebJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotWebJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Web Job ID
func (id SlotWebJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/webJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.WebJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Web Job ID
func (id SlotWebJobId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticWebJobs", "webJobs", "webJobs"),
		resourceids.UserSpecifiedSegment("webJobName", "webJobName"),
	}
}

// String returns a human-readable description of this Slot Web Job ID
func (id SlotWebJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Web Job Name: %q", id.WebJobName),
	}
	return fmt.Sprintf("Slot Web Job (%s)", strings.Join(components, "\n"))
}
