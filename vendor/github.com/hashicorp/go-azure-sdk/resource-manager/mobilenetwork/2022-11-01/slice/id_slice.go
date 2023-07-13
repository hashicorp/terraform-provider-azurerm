package slice

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SliceId{}

// SliceId is a struct representing the Resource ID for a Slice
type SliceId struct {
	SubscriptionId    string
	ResourceGroupName string
	MobileNetworkName string
	SliceName         string
}

// NewSliceID returns a new SliceId struct
func NewSliceID(subscriptionId string, resourceGroupName string, mobileNetworkName string, sliceName string) SliceId {
	return SliceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MobileNetworkName: mobileNetworkName,
		SliceName:         sliceName,
	}
}

// ParseSliceID parses 'input' into a SliceId
func ParseSliceID(input string) (*SliceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SliceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SliceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.SliceName, ok = parsed.Parsed["sliceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sliceName", *parsed)
	}

	return &id, nil
}

// ParseSliceIDInsensitively parses 'input' case-insensitively into a SliceId
// note: this method should only be used for API response data and not user input
func ParseSliceIDInsensitively(input string) (*SliceId, error) {
	parser := resourceids.NewParserFromResourceIdType(SliceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SliceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.SliceName, ok = parsed.Parsed["sliceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sliceName", *parsed)
	}

	return &id, nil
}

// ValidateSliceID checks that 'input' can be parsed as a Slice ID
func ValidateSliceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSliceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slice ID
func (id SliceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/mobileNetworks/%s/slices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName, id.SliceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slice ID
func (id SliceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticMobileNetworks", "mobileNetworks", "mobileNetworks"),
		resourceids.UserSpecifiedSegment("mobileNetworkName", "mobileNetworkValue"),
		resourceids.StaticSegment("staticSlices", "slices", "slices"),
		resourceids.UserSpecifiedSegment("sliceName", "sliceValue"),
	}
}

// String returns a human-readable description of this Slice ID
func (id SliceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Mobile Network Name: %q", id.MobileNetworkName),
		fmt.Sprintf("Slice Name: %q", id.SliceName),
	}
	return fmt.Sprintf("Slice (%s)", strings.Join(components, "\n"))
}
