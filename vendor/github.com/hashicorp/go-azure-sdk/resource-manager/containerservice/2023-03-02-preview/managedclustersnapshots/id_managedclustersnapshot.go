package managedclustersnapshots

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedClusterSnapshotId{})
}

var _ resourceids.ResourceId = &ManagedClusterSnapshotId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ManagedClusterSnapshotId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedClusterSnapshotId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedClusterSnapshotIDInsensitively parses 'input' case-insensitively into a ManagedClusterSnapshotId
// note: this method should only be used for API response data and not user input
func ParseManagedClusterSnapshotIDInsensitively(input string) (*ManagedClusterSnapshotId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedClusterSnapshotId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedClusterSnapshotId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedClusterSnapshotId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedClusterSnapshotName, ok = input.Parsed["managedClusterSnapshotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedClusterSnapshotName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("managedClusterSnapshotName", "managedClusterSnapshotName"),
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
