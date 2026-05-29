package marketplacegalleryimages

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MarketplaceGalleryImageId{})
}

var _ resourceids.ResourceId = &MarketplaceGalleryImageId{}

// MarketplaceGalleryImageId is a struct representing the Resource ID for a Marketplace Gallery Image
type MarketplaceGalleryImageId struct {
	SubscriptionId              string
	ResourceGroupName           string
	MarketplaceGalleryImageName string
}

// NewMarketplaceGalleryImageID returns a new MarketplaceGalleryImageId struct
func NewMarketplaceGalleryImageID(subscriptionId string, resourceGroupName string, marketplaceGalleryImageName string) MarketplaceGalleryImageId {
	return MarketplaceGalleryImageId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		MarketplaceGalleryImageName: marketplaceGalleryImageName,
	}
}

// ParseMarketplaceGalleryImageID parses 'input' into a MarketplaceGalleryImageId
func ParseMarketplaceGalleryImageID(input string) (*MarketplaceGalleryImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MarketplaceGalleryImageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MarketplaceGalleryImageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMarketplaceGalleryImageIDInsensitively parses 'input' case-insensitively into a MarketplaceGalleryImageId
// note: this method should only be used for API response data and not user input
func ParseMarketplaceGalleryImageIDInsensitively(input string) (*MarketplaceGalleryImageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MarketplaceGalleryImageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MarketplaceGalleryImageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MarketplaceGalleryImageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MarketplaceGalleryImageName, ok = input.Parsed["marketplaceGalleryImageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "marketplaceGalleryImageName", input)
	}

	return nil
}

// ValidateMarketplaceGalleryImageID checks that 'input' can be parsed as a Marketplace Gallery Image ID
func ValidateMarketplaceGalleryImageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMarketplaceGalleryImageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Marketplace Gallery Image ID
func (id MarketplaceGalleryImageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/marketplaceGalleryImages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MarketplaceGalleryImageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Marketplace Gallery Image ID
func (id MarketplaceGalleryImageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticMarketplaceGalleryImages", "marketplaceGalleryImages", "marketplaceGalleryImages"),
		resourceids.UserSpecifiedSegment("marketplaceGalleryImageName", "marketplaceGalleryImageName"),
	}
}

// String returns a human-readable description of this Marketplace Gallery Image ID
func (id MarketplaceGalleryImageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Marketplace Gallery Image Name: %q", id.MarketplaceGalleryImageName),
	}
	return fmt.Sprintf("Marketplace Gallery Image (%s)", strings.Join(components, "\n"))
}
