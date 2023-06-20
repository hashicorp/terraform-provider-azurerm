package snapshots

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SnapshotId{}

// SnapshotId is a struct representing the Resource ID for a Snapshot
type SnapshotId struct {
	SubscriptionId    string
	ResourceGroupName string
	SnapshotName      string
}

// NewSnapshotID returns a new SnapshotId struct
func NewSnapshotID(subscriptionId string, resourceGroupName string, snapshotName string) SnapshotId {
	return SnapshotId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SnapshotName:      snapshotName,
	}
}

// ParseSnapshotID parses 'input' into a SnapshotId
func ParseSnapshotID(input string) (*SnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(SnapshotId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SnapshotId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SnapshotName, ok = parsed.Parsed["snapshotName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "snapshotName", *parsed)
	}

	return &id, nil
}

// ParseSnapshotIDInsensitively parses 'input' case-insensitively into a SnapshotId
// note: this method should only be used for API response data and not user input
func ParseSnapshotIDInsensitively(input string) (*SnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(SnapshotId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SnapshotId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SnapshotName, ok = parsed.Parsed["snapshotName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "snapshotName", *parsed)
	}

	return &id, nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SnapshotName)
}

// Segments returns a slice of Resource ID Segments which comprise this Snapshot ID
func (id SnapshotId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticSnapshots", "snapshots", "snapshots"),
		resourceids.UserSpecifiedSegment("snapshotName", "snapshotValue"),
	}
}

// String returns a human-readable description of this Snapshot ID
func (id SnapshotId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Snapshot Name: %q", id.SnapshotName),
	}
	return fmt.Sprintf("Snapshot (%s)", strings.Join(components, "\n"))
}
