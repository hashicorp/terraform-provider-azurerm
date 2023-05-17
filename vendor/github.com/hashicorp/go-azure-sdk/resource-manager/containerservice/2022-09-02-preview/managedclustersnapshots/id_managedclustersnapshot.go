package managedclustersnapshots

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ManagedClusterSnapshotId{}

// ManagedClusterSnapshotId is a struct representing the Resource ID for a Managed Cluster Snapshot
type ManagedClusterSnapshotId struct {
	SubscriptionId             string
	ResourceGroupName          string
	ManagedClusterSnapshotName string
}

// NewManagedClusterSnapshotID returns a new ManagedClusterSnapshotId struct
func NewManagedClusterSnapshotID(subscriptionId string, resourceGroupName string, managedClusterSnapshotName string) ManagedClusterSnapshotId {
	return ManagedClusterSnapshotId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		ManagedClusterSnapshotName: managedClusterSnapshotName,
	}
}

// ParseManagedClusterSnapshotID parses 'input' into a ManagedClusterSnapshotId
func ParseManagedClusterSnapshotID(input string) (*ManagedClusterSnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedClusterSnapshotId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedClusterSnapshotId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterSnapshotName, ok = parsed.Parsed["managedClusterSnapshotName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterSnapshotName", *parsed)
	}

	return &id, nil
}

// ParseManagedClusterSnapshotIDInsensitively parses 'input' case-insensitively into a ManagedClusterSnapshotId
// note: this method should only be used for API response data and not user input
func ParseManagedClusterSnapshotIDInsensitively(input string) (*ManagedClusterSnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedClusterSnapshotId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedClusterSnapshotId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterSnapshotName, ok = parsed.Parsed["managedClusterSnapshotName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterSnapshotName", *parsed)
	}

	return &id, nil
}

// ValidateManagedClusterSnapshotID checks that 'input' can be parsed as a Managed Cluster Snapshot ID
func ValidateManagedClusterSnapshotID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedClusterSnapshotID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Cluster Snapshot ID
func (id ManagedClusterSnapshotId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusterSnapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterSnapshotName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Cluster Snapshot ID
func (id ManagedClusterSnapshotId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusterSnapshots", "managedClusterSnapshots", "managedClusterSnapshots"),
		resourceids.UserSpecifiedSegment("managedClusterSnapshotName", "managedClusterSnapshotValue"),
	}
}

// String returns a human-readable description of this Managed Cluster Snapshot ID
func (id ManagedClusterSnapshotId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Snapshot Name: %q", id.ManagedClusterSnapshotName),
	}
	return fmt.Sprintf("Managed Cluster Snapshot (%s)", strings.Join(components, "\n"))
}
