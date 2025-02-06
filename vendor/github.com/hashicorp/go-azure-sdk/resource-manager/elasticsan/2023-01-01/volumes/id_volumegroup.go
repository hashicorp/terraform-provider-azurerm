package volumes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VolumeGroupId{})
}

var _ resourceids.ResourceId = &VolumeGroupId{}

// VolumeGroupId is a struct representing the Resource ID for a Volume Group
type VolumeGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	ElasticSanName    string
	VolumeGroupName   string
}

// NewVolumeGroupID returns a new VolumeGroupId struct
func NewVolumeGroupID(subscriptionId string, resourceGroupName string, elasticSanName string, volumeGroupName string) VolumeGroupId {
	return VolumeGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ElasticSanName:    elasticSanName,
		VolumeGroupName:   volumeGroupName,
	}
}

// ParseVolumeGroupID parses 'input' into a VolumeGroupId
func ParseVolumeGroupID(input string) (*VolumeGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VolumeGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VolumeGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVolumeGroupIDInsensitively parses 'input' case-insensitively into a VolumeGroupId
// note: this method should only be used for API response data and not user input
func ParseVolumeGroupIDInsensitively(input string) (*VolumeGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VolumeGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VolumeGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VolumeGroupId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateVolumeGroupID checks that 'input' can be parsed as a Volume Group ID
func ValidateVolumeGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVolumeGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Volume Group ID
func (id VolumeGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ElasticSan/elasticSans/%s/volumeGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName, id.VolumeGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Volume Group ID
func (id VolumeGroupId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Volume Group ID
func (id VolumeGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Elastic San Name: %q", id.ElasticSanName),
		fmt.Sprintf("Volume Group Name: %q", id.VolumeGroupName),
	}
	return fmt.Sprintf("Volume Group (%s)", strings.Join(components, "\n"))
}
