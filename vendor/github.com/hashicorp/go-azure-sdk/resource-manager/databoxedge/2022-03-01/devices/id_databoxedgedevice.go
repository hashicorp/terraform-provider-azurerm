package devices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DataBoxEdgeDeviceId{}

// DataBoxEdgeDeviceId is a struct representing the Resource ID for a Data Box Edge Device
type DataBoxEdgeDeviceId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DataBoxEdgeDeviceName string
}

// NewDataBoxEdgeDeviceID returns a new DataBoxEdgeDeviceId struct
func NewDataBoxEdgeDeviceID(subscriptionId string, resourceGroupName string, dataBoxEdgeDeviceName string) DataBoxEdgeDeviceId {
	return DataBoxEdgeDeviceId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DataBoxEdgeDeviceName: dataBoxEdgeDeviceName,
	}
}

// ParseDataBoxEdgeDeviceID parses 'input' into a DataBoxEdgeDeviceId
func ParseDataBoxEdgeDeviceID(input string) (*DataBoxEdgeDeviceId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataBoxEdgeDeviceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataBoxEdgeDeviceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DataBoxEdgeDeviceName, ok = parsed.Parsed["dataBoxEdgeDeviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataBoxEdgeDeviceName", *parsed)
	}

	return &id, nil
}

// ParseDataBoxEdgeDeviceIDInsensitively parses 'input' case-insensitively into a DataBoxEdgeDeviceId
// note: this method should only be used for API response data and not user input
func ParseDataBoxEdgeDeviceIDInsensitively(input string) (*DataBoxEdgeDeviceId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataBoxEdgeDeviceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataBoxEdgeDeviceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DataBoxEdgeDeviceName, ok = parsed.Parsed["dataBoxEdgeDeviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataBoxEdgeDeviceName", *parsed)
	}

	return &id, nil
}

// ValidateDataBoxEdgeDeviceID checks that 'input' can be parsed as a Data Box Edge Device ID
func ValidateDataBoxEdgeDeviceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataBoxEdgeDeviceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Box Edge Device ID
func (id DataBoxEdgeDeviceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataBoxEdge/dataBoxEdgeDevices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DataBoxEdgeDeviceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Box Edge Device ID
func (id DataBoxEdgeDeviceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataBoxEdge", "Microsoft.DataBoxEdge", "Microsoft.DataBoxEdge"),
		resourceids.StaticSegment("staticDataBoxEdgeDevices", "dataBoxEdgeDevices", "dataBoxEdgeDevices"),
		resourceids.UserSpecifiedSegment("dataBoxEdgeDeviceName", "dataBoxEdgeDeviceValue"),
	}
}

// String returns a human-readable description of this Data Box Edge Device ID
func (id DataBoxEdgeDeviceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Data Box Edge Device Name: %q", id.DataBoxEdgeDeviceName),
	}
	return fmt.Sprintf("Data Box Edge Device (%s)", strings.Join(components, "\n"))
}
