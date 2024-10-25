package resource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SpatialAnchorsAccountId{})
}

var _ resourceids.ResourceId = &SpatialAnchorsAccountId{}

// SpatialAnchorsAccountId is a struct representing the Resource ID for a Spatial Anchors Account
type SpatialAnchorsAccountId struct {
	SubscriptionId            string
	ResourceGroupName         string
	SpatialAnchorsAccountName string
}

// NewSpatialAnchorsAccountID returns a new SpatialAnchorsAccountId struct
func NewSpatialAnchorsAccountID(subscriptionId string, resourceGroupName string, spatialAnchorsAccountName string) SpatialAnchorsAccountId {
	return SpatialAnchorsAccountId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		SpatialAnchorsAccountName: spatialAnchorsAccountName,
	}
}

// ParseSpatialAnchorsAccountID parses 'input' into a SpatialAnchorsAccountId
func ParseSpatialAnchorsAccountID(input string) (*SpatialAnchorsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SpatialAnchorsAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpatialAnchorsAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSpatialAnchorsAccountIDInsensitively parses 'input' case-insensitively into a SpatialAnchorsAccountId
// note: this method should only be used for API response data and not user input
func ParseSpatialAnchorsAccountIDInsensitively(input string) (*SpatialAnchorsAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SpatialAnchorsAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpatialAnchorsAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SpatialAnchorsAccountId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpatialAnchorsAccountName, ok = input.Parsed["spatialAnchorsAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "spatialAnchorsAccountName", input)
	}

	return nil
}

// ValidateSpatialAnchorsAccountID checks that 'input' can be parsed as a Spatial Anchors Account ID
func ValidateSpatialAnchorsAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSpatialAnchorsAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Spatial Anchors Account ID
func (id SpatialAnchorsAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpatialAnchorsAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Spatial Anchors Account ID
func (id SpatialAnchorsAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMixedReality", "Microsoft.MixedReality", "Microsoft.MixedReality"),
		resourceids.StaticSegment("staticSpatialAnchorsAccounts", "spatialAnchorsAccounts", "spatialAnchorsAccounts"),
		resourceids.UserSpecifiedSegment("spatialAnchorsAccountName", "spatialAnchorsAccountName"),
	}
}

// String returns a human-readable description of this Spatial Anchors Account ID
func (id SpatialAnchorsAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spatial Anchors Account Name: %q", id.SpatialAnchorsAccountName),
	}
	return fmt.Sprintf("Spatial Anchors Account (%s)", strings.Join(components, "\n"))
}
