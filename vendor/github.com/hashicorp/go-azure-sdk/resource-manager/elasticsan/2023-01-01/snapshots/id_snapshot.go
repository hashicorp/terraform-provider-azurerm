package snapshots

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
	ElasticSanName    string
	VolumeGroupName   string
	SnapshotName      string
}

// NewSnapshotID returns a new SnapshotId struct
func NewSnapshotID(subscriptionId string, resourceGroupName string, elasticSanName string, volumeGroupName string, snapshotName string) SnapshotId {
	return SnapshotId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ElasticSanName:    elasticSanName,
		VolumeGroupName:   volumeGroupName,
		SnapshotName:      snapshotName,
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

	if id.ElasticSanName, ok = input.Parsed["elasticSanName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "elasticSanName", input)
	}

	if id.VolumeGroupName, ok = input.Parsed["volumeGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "volumeGroupName", input)
	}

	if id.SnapshotName, ok = input.Parsed["snapshotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "snapshotName", input)
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ElasticSan/elasticSans/%s/volumeGroups/%s/snapshots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName, id.VolumeGroupName, id.SnapshotName)
}

// Segments returns a slice of Resource ID Segments which comprise this Snapshot ID
func (id SnapshotId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftElasticSan", "Microsoft.ElasticSan", "Microsoft.ElasticSan"),
		resourceids.StaticSegment("staticElasticSans", "elasticSans", "elasticSans"),
		resourceids.UserSpecifiedSegment("elasticSanName", "elasticSanName"),
		resourceids.StaticSegment("staticVolumeGroups", "volumeGroups", "volumeGroups"),
		resourceids.UserSpecifiedSegment("volumeGroupName", "volumeGroupName"),
		resourceids.StaticSegment("staticSnapshots", "snapshots", "snapshots"),
		resourceids.UserSpecifiedSegment("snapshotName", "snapshotName"),
	}
}

// String returns a human-readable description of this Snapshot ID
func (id SnapshotId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Elastic San Name: %q", id.ElasticSanName),
		fmt.Sprintf("Volume Group Name: %q", id.VolumeGroupName),
		fmt.Sprintf("Snapshot Name: %q", id.SnapshotName),
	}
	return fmt.Sprintf("Snapshot (%s)", strings.Join(components, "\n"))
}
