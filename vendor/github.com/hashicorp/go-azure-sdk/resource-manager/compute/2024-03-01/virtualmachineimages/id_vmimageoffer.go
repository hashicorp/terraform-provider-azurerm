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
	recaser.RegisterResourceId(&VMImageOfferId{})
}

var _ resourceids.ResourceId = &VMImageOfferId{}

// VMImageOfferId is a struct representing the Resource ID for a V M Image Offer
type VMImageOfferId struct {
	SubscriptionId string
	LocationName   string
	EdgeZoneName   string
	PublisherName  string
	OfferName      string
}

// NewVMImageOfferID returns a new VMImageOfferId struct
func NewVMImageOfferID(subscriptionId string, locationName string, edgeZoneName string, publisherName string, offerName string) VMImageOfferId {
	return VMImageOfferId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		EdgeZoneName:   edgeZoneName,
		PublisherName:  publisherName,
		OfferName:      offerName,
	}
}

// ParseVMImageOfferID parses 'input' into a VMImageOfferId
func ParseVMImageOfferID(input string) (*VMImageOfferId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMImageOfferId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMImageOfferId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVMImageOfferIDInsensitively parses 'input' case-insensitively into a VMImageOfferId
// note: this method should only be used for API response data and not user input
func ParseVMImageOfferIDInsensitively(input string) (*VMImageOfferId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VMImageOfferId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VMImageOfferId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VMImageOfferId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateVMImageOfferID checks that 'input' can be parsed as a V M Image Offer ID
func ValidateVMImageOfferID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVMImageOfferID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted V M Image Offer ID
func (id VMImageOfferId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Compute/locations/%s/edgeZones/%s/publishers/%s/artifactTypes/vmImage/offers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.EdgeZoneName, id.PublisherName, id.OfferName)
}

// Segments returns a slice of Resource ID Segments which comprise this V M Image Offer ID
func (id VMImageOfferId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this V M Image Offer ID
func (id VMImageOfferId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Edge Zone Name: %q", id.EdgeZoneName),
		fmt.Sprintf("Publisher Name: %q", id.PublisherName),
		fmt.Sprintf("Offer Name: %q", id.OfferName),
	}
	return fmt.Sprintf("V M Image Offer (%s)", strings.Join(components, "\n"))
}
