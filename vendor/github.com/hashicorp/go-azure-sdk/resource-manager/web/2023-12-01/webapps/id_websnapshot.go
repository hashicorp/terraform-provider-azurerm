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
	recaser.RegisterResourceId(&WebSnapshotId{})
}

var _ resourceids.ResourceId = &WebSnapshotId{}

// WebSnapshotId is a struct representing the Resource ID for a Web Snapshot
type WebSnapshotId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	SnapshotId        string
}

// NewWebSnapshotID returns a new WebSnapshotId struct
func NewWebSnapshotID(subscriptionId string, resourceGroupName string, siteName string, slotName string, snapshotId string) WebSnapshotId {
	return WebSnapshotId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		SnapshotId:        snapshotId,
	}
}

// ParseWebSnapshotID parses 'input' into a WebSnapshotId
func ParseWebSnapshotID(input string) (*WebSnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WebSnapshotId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WebSnapshotId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWebSnapshotIDInsensitively parses 'input' case-insensitively into a WebSnapshotId
// note: this method should only be used for API response data and not user input
func ParseWebSnapshotIDInsensitively(input string) (*WebSnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WebSnapshotId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WebSnapshotId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WebSnapshotId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SnapshotId, ok = input.Parsed["snapshotId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "snapshotId", input)
	}

	return nil
}

// ValidateWebSnapshotID checks that 'input' can be parsed as a Web Snapshot ID
func ValidateWebSnapshotID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWebSnapshotID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Web Snapshot ID
func (id WebSnapshotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/config/web/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.SnapshotId)
}

// Segments returns a slice of Resource ID Segments which comprise this Web Snapshot ID
func (id WebSnapshotId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticConfig", "config", "config"),
		resourceids.StaticSegment("staticWeb", "web", "web"),
		resourceids.StaticSegment("staticSnapshots", "snapshots", "snapshots"),
		resourceids.UserSpecifiedSegment("snapshotId", "snapshotId"),
	}
}

// String returns a human-readable description of this Web Snapshot ID
func (id WebSnapshotId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Snapshot: %q", id.SnapshotId),
	}
	return fmt.Sprintf("Web Snapshot (%s)", strings.Join(components, "\n"))
}
