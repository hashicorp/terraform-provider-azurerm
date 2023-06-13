package datanetwork

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DataNetworkId{}

// DataNetworkId is a struct representing the Resource ID for a Data Network
type DataNetworkId struct {
	SubscriptionId    string
	ResourceGroupName string
	MobileNetworkName string
	DataNetworkName   string
}

// NewDataNetworkID returns a new DataNetworkId struct
func NewDataNetworkID(subscriptionId string, resourceGroupName string, mobileNetworkName string, dataNetworkName string) DataNetworkId {
	return DataNetworkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MobileNetworkName: mobileNetworkName,
		DataNetworkName:   dataNetworkName,
	}
}

// ParseDataNetworkID parses 'input' into a DataNetworkId
func ParseDataNetworkID(input string) (*DataNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataNetworkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.DataNetworkName, ok = parsed.Parsed["dataNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataNetworkName", *parsed)
	}

	return &id, nil
}

// ParseDataNetworkIDInsensitively parses 'input' case-insensitively into a DataNetworkId
// note: this method should only be used for API response data and not user input
func ParseDataNetworkIDInsensitively(input string) (*DataNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataNetworkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.DataNetworkName, ok = parsed.Parsed["dataNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataNetworkName", *parsed)
	}

	return &id, nil
}

// ValidateDataNetworkID checks that 'input' can be parsed as a Data Network ID
func ValidateDataNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Network ID
func (id DataNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/mobileNetworks/%s/dataNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName, id.DataNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Network ID
func (id DataNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticMobileNetworks", "mobileNetworks", "mobileNetworks"),
		resourceids.UserSpecifiedSegment("mobileNetworkName", "mobileNetworkValue"),
		resourceids.StaticSegment("staticDataNetworks", "dataNetworks", "dataNetworks"),
		resourceids.UserSpecifiedSegment("dataNetworkName", "dataNetworkValue"),
	}
}

// String returns a human-readable description of this Data Network ID
func (id DataNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Mobile Network Name: %q", id.MobileNetworkName),
		fmt.Sprintf("Data Network Name: %q", id.DataNetworkName),
	}
	return fmt.Sprintf("Data Network (%s)", strings.Join(components, "\n"))
}
