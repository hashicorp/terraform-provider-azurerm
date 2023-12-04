package volumes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VolumeId{}

// VolumeId is a struct representing the Resource ID for a Volume
type VolumeId struct {
	SubscriptionId    string
	ResourceGroupName string
	ElasticSanName    string
	VolumeGroupName   string
	VolumeName        string
}

// NewVolumeID returns a new VolumeId struct
func NewVolumeID(subscriptionId string, resourceGroupName string, elasticSanName string, volumeGroupName string, volumeName string) VolumeId {
	return VolumeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ElasticSanName:    elasticSanName,
		VolumeGroupName:   volumeGroupName,
		VolumeName:        volumeName,
	}
}

// ParseVolumeID parses 'input' into a VolumeId
func ParseVolumeID(input string) (*VolumeId, error) {
	parser := resourceids.NewParserFromResourceIdType(VolumeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VolumeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ElasticSanName, ok = parsed.Parsed["elasticSanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "elasticSanName", *parsed)
	}

	if id.VolumeGroupName, ok = parsed.Parsed["volumeGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "volumeGroupName", *parsed)
	}

	if id.VolumeName, ok = parsed.Parsed["volumeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "volumeName", *parsed)
	}

	return &id, nil
}

// ParseVolumeIDInsensitively parses 'input' case-insensitively into a VolumeId
// note: this method should only be used for API response data and not user input
func ParseVolumeIDInsensitively(input string) (*VolumeId, error) {
	parser := resourceids.NewParserFromResourceIdType(VolumeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VolumeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ElasticSanName, ok = parsed.Parsed["elasticSanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "elasticSanName", *parsed)
	}

	if id.VolumeGroupName, ok = parsed.Parsed["volumeGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "volumeGroupName", *parsed)
	}

	if id.VolumeName, ok = parsed.Parsed["volumeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "volumeName", *parsed)
	}

	return &id, nil
}

// ValidateVolumeID checks that 'input' can be parsed as a Volume ID
func ValidateVolumeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVolumeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Volume ID
func (id VolumeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ElasticSan/elasticSans/%s/volumeGroups/%s/volumes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName, id.VolumeGroupName, id.VolumeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Volume ID
func (id VolumeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftElasticSan", "Microsoft.ElasticSan", "Microsoft.ElasticSan"),
		resourceids.StaticSegment("staticElasticSans", "elasticSans", "elasticSans"),
		resourceids.UserSpecifiedSegment("elasticSanName", "elasticSanValue"),
		resourceids.StaticSegment("staticVolumeGroups", "volumeGroups", "volumeGroups"),
		resourceids.UserSpecifiedSegment("volumeGroupName", "volumeGroupValue"),
		resourceids.StaticSegment("staticVolumes", "volumes", "volumes"),
		resourceids.UserSpecifiedSegment("volumeName", "volumeValue"),
	}
}

// String returns a human-readable description of this Volume ID
func (id VolumeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Elastic San Name: %q", id.ElasticSanName),
		fmt.Sprintf("Volume Group Name: %q", id.VolumeGroupName),
		fmt.Sprintf("Volume Name: %q", id.VolumeName),
	}
	return fmt.Sprintf("Volume (%s)", strings.Join(components, "\n"))
}
