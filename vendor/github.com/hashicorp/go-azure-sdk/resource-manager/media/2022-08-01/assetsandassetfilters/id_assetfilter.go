package assetsandassetfilters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AssetFilterId{}

// AssetFilterId is a struct representing the Resource ID for a Asset Filter
type AssetFilterId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	AssetName         string
	AssetFilterName   string
}

// NewAssetFilterID returns a new AssetFilterId struct
func NewAssetFilterID(subscriptionId string, resourceGroupName string, mediaServiceName string, assetName string, assetFilterName string) AssetFilterId {
	return AssetFilterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		AssetName:         assetName,
		AssetFilterName:   assetFilterName,
	}
}

// ParseAssetFilterID parses 'input' into a AssetFilterId
func ParseAssetFilterID(input string) (*AssetFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssetFilterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssetFilterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.AssetName, ok = parsed.Parsed["assetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assetName", *parsed)
	}

	if id.AssetFilterName, ok = parsed.Parsed["assetFilterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assetFilterName", *parsed)
	}

	return &id, nil
}

// ParseAssetFilterIDInsensitively parses 'input' case-insensitively into a AssetFilterId
// note: this method should only be used for API response data and not user input
func ParseAssetFilterIDInsensitively(input string) (*AssetFilterId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssetFilterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssetFilterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.AssetName, ok = parsed.Parsed["assetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assetName", *parsed)
	}

	if id.AssetFilterName, ok = parsed.Parsed["assetFilterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "assetFilterName", *parsed)
	}

	return &id, nil
}

// ValidateAssetFilterID checks that 'input' can be parsed as a Asset Filter ID
func ValidateAssetFilterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssetFilterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Asset Filter ID
func (id AssetFilterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/assets/%s/assetFilters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.AssetName, id.AssetFilterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Asset Filter ID
func (id AssetFilterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticAssets", "assets", "assets"),
		resourceids.UserSpecifiedSegment("assetName", "assetValue"),
		resourceids.StaticSegment("staticAssetFilters", "assetFilters", "assetFilters"),
		resourceids.UserSpecifiedSegment("assetFilterName", "assetFilterValue"),
	}
}

// String returns a human-readable description of this Asset Filter ID
func (id AssetFilterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Asset Name: %q", id.AssetName),
		fmt.Sprintf("Asset Filter Name: %q", id.AssetFilterName),
	}
	return fmt.Sprintf("Asset Filter (%s)", strings.Join(components, "\n"))
}
