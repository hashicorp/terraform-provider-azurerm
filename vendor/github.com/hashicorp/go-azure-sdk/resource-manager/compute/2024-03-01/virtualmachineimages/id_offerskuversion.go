package virtualmachineimages

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OfferSkuVersionId{})
}

var _ resourceids.ResourceId = &OfferSkuVersionId{}

// OfferSkuVersionId is a struct representing the Resource ID for a Offer Sku Version
type OfferSkuVersionId struct {
	SubscriptionId string
	LocationName   string
	EdgeZoneName   string
	PublisherName  string
	OfferName      string
	SkuName        string
	VersionName    string
}

// NewOfferSkuVersionID returns a new OfferSkuVersionId struct
func NewOfferSkuVersionID(subscriptionId string, locationName string, edgeZoneName string, publisherName string, offerName string, skuName string, versionName string) OfferSkuVersionId {
	return OfferSkuVersionId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		EdgeZoneName:   edgeZoneName,
		PublisherName:  publisherName,
		OfferName:      offerName,
		SkuName:        skuName,
		VersionName:    versionName,
	}
}

// ParseOfferSkuVersionID parses 'input' into a OfferSkuVersionId
func ParseOfferSkuVersionID(input string) (*OfferSkuVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OfferSkuVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OfferSkuVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOfferSkuVersionIDInsensitively parses 'input' case-insensitively into a OfferSkuVersionId
// note: this method should only be used for API response data and not user input
func ParseOfferSkuVersionIDInsensitively(input string) (*OfferSkuVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OfferSkuVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OfferSkuVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OfferSkuVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.EdgeZoneName, ok = input.Parsed["edgeZoneName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "edgeZoneName", input)
	}

	if id.PublisherName, ok = input.Parsed["publisherName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "publisherName", input)
	}

	if id.OfferName, ok = input.Parsed["offerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "offerName", input)
	}

	if id.SkuName, ok = input.Parsed["skuName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "skuName", input)
	}

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
}

// ValidateOfferSkuVersionID checks that 'input' can be parsed as a Offer Sku Version ID
func ValidateOfferSkuVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOfferSkuVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Offer Sku Version ID
func (id OfferSkuVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Compute/locations/%s/edgeZones/%s/publishers/%s/artifactTypes/vmImage/offers/%s/skus/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.EdgeZoneName, id.PublisherName, id.OfferName, id.SkuName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Offer Sku Version ID
func (id OfferSkuVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticEdgeZones", "edgeZones", "edgeZones"),
		resourceids.UserSpecifiedSegment("edgeZoneName", "edgeZoneName"),
		resourceids.StaticSegment("staticPublishers", "publishers", "publishers"),
		resourceids.UserSpecifiedSegment("publisherName", "publisherName"),
		resourceids.StaticSegment("staticArtifactTypes", "artifactTypes", "artifactTypes"),
		resourceids.StaticSegment("staticVmImage", "vmImage", "vmImage"),
		resourceids.StaticSegment("staticOffers", "offers", "offers"),
		resourceids.UserSpecifiedSegment("offerName", "offerName"),
		resourceids.StaticSegment("staticSkus", "skus", "skus"),
		resourceids.UserSpecifiedSegment("skuName", "skuName"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionName"),
	}
}

// String returns a human-readable description of this Offer Sku Version ID
func (id OfferSkuVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Edge Zone Name: %q", id.EdgeZoneName),
		fmt.Sprintf("Publisher Name: %q", id.PublisherName),
		fmt.Sprintf("Offer Name: %q", id.OfferName),
		fmt.Sprintf("Sku Name: %q", id.SkuName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Offer Sku Version (%s)", strings.Join(components, "\n"))
}
