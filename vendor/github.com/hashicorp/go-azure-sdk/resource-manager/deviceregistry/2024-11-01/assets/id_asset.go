package assets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AssetId{})
}

var _ resourceids.ResourceId = &AssetId{}

// AssetId is a struct representing the Resource ID for a Asset
type AssetId struct {
	SubscriptionId    string
	ResourceGroupName string
	AssetName         string
}

// NewAssetID returns a new AssetId struct
func NewAssetID(subscriptionId string, resourceGroupName string, assetName string) AssetId {
	return AssetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AssetName:         assetName,
	}
}

// ParseAssetID parses 'input' into a AssetId
func ParseAssetID(input string) (*AssetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAssetIDInsensitively parses 'input' case-insensitively into a AssetId
// note: this method should only be used for API response data and not user input
func ParseAssetIDInsensitively(input string) (*AssetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AssetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AssetName, ok = input.Parsed["assetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "assetName", input)
	}

	return nil
}

// ValidateAssetID checks that 'input' can be parsed as a Asset ID
func ValidateAssetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Asset ID
func (id AssetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DeviceRegistry/assets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AssetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Asset ID
func (id AssetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDeviceRegistry", "Microsoft.DeviceRegistry", "Microsoft.DeviceRegistry"),
		resourceids.StaticSegment("staticAssets", "assets", "assets"),
		resourceids.UserSpecifiedSegment("assetName", "assetName"),
	}
}

// String returns a human-readable description of this Asset ID
func (id AssetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Asset Name: %q", id.AssetName),
	}
	return fmt.Sprintf("Asset (%s)", strings.Join(components, "\n"))
}
