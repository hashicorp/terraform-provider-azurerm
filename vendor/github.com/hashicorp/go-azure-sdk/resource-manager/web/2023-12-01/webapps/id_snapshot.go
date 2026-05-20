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
	recaser.RegisterResourceId(&SnapshotId{})
}

var _ resourceids.ResourceId = &SnapshotId{}

// SnapshotId is a struct representing the Resource ID for a Snapshot
type SnapshotId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SnapshotId        string
}

// NewSnapshotID returns a new SnapshotId struct
func NewSnapshotID(subscriptionId string, resourceGroupName string, siteName string, snapshotId string) SnapshotId {
	return SnapshotId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SnapshotId:        snapshotId,
	}
}

// ParseSnapshotID parses 'input' into a SnapshotId
func ParseSnapshotID(input string) (*SnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SnapshotId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SnapshotId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSnapshotIDInsensitively parses 'input' case-insensitively into a SnapshotId
// note: this method should only be used for API response data and not user input
func ParseSnapshotIDInsensitively(input string) (*SnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SnapshotId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SnapshotId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SnapshotId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SnapshotId, ok = input.Parsed["snapshotId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "snapshotId", input)
	}

	return nil
}

// ValidateSnapshotID checks that 'input' can be parsed as a Snapshot ID
func ValidateSnapshotID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSnapshotID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Snapshot ID
func (id SnapshotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/config/web/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SnapshotId)
}

// Segments returns a slice of Resource ID Segments which comprise this Snapshot ID
func (id SnapshotId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticConfig", "config", "config"),
		resourceids.StaticSegment("staticWeb", "web", "web"),
		resourceids.StaticSegment("staticSnapshots", "snapshots", "snapshots"),
		resourceids.UserSpecifiedSegment("snapshotId", "snapshotId"),
	}
}

// String returns a human-readable description of this Snapshot ID
func (id SnapshotId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Snapshot: %q", id.SnapshotId),
	}
	return fmt.Sprintf("Snapshot (%s)", strings.Join(components, "\n"))
}
