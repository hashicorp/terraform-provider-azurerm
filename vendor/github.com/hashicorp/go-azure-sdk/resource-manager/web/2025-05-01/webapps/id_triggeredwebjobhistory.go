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
	recaser.RegisterResourceId(&TriggeredWebJobHistoryId{})
}

var _ resourceids.ResourceId = &TriggeredWebJobHistoryId{}

// TriggeredWebJobHistoryId is a struct representing the Resource ID for a Triggered Web Job History
type TriggeredWebJobHistoryId struct {
	SubscriptionId      string
	ResourceGroupName   string
	SiteName            string
	SlotName            string
	TriggeredWebJobName string
	HistoryName         string
}

// NewTriggeredWebJobHistoryID returns a new TriggeredWebJobHistoryId struct
func NewTriggeredWebJobHistoryID(subscriptionId string, resourceGroupName string, siteName string, slotName string, triggeredWebJobName string, historyName string) TriggeredWebJobHistoryId {
	return TriggeredWebJobHistoryId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		SiteName:            siteName,
		SlotName:            slotName,
		TriggeredWebJobName: triggeredWebJobName,
		HistoryName:         historyName,
	}
}

// ParseTriggeredWebJobHistoryID parses 'input' into a TriggeredWebJobHistoryId
func ParseTriggeredWebJobHistoryID(input string) (*TriggeredWebJobHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TriggeredWebJobHistoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TriggeredWebJobHistoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTriggeredWebJobHistoryIDInsensitively parses 'input' case-insensitively into a TriggeredWebJobHistoryId
// note: this method should only be used for API response data and not user input
func ParseTriggeredWebJobHistoryIDInsensitively(input string) (*TriggeredWebJobHistoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TriggeredWebJobHistoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TriggeredWebJobHistoryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TriggeredWebJobHistoryId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TriggeredWebJobName, ok = input.Parsed["triggeredWebJobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "triggeredWebJobName", input)
	}

	if id.HistoryName, ok = input.Parsed["historyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "historyName", input)
	}

	return nil
}

// ValidateTriggeredWebJobHistoryID checks that 'input' can be parsed as a Triggered Web Job History ID
func ValidateTriggeredWebJobHistoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTriggeredWebJobHistoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Triggered Web Job History ID
func (id TriggeredWebJobHistoryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/triggeredWebJobs/%s/history/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.TriggeredWebJobName, id.HistoryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Triggered Web Job History ID
func (id TriggeredWebJobHistoryId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticTriggeredWebJobs", "triggeredWebJobs", "triggeredWebJobs"),
		resourceids.UserSpecifiedSegment("triggeredWebJobName", "triggeredWebJobName"),
		resourceids.StaticSegment("staticHistory", "history", "history"),
		resourceids.UserSpecifiedSegment("historyName", "historyName"),
	}
}

// String returns a human-readable description of this Triggered Web Job History ID
func (id TriggeredWebJobHistoryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Triggered Web Job Name: %q", id.TriggeredWebJobName),
		fmt.Sprintf("History Name: %q", id.HistoryName),
	}
	return fmt.Sprintf("Triggered Web Job History (%s)", strings.Join(components, "\n"))
}
