package ipallocations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = IPAllocationId{}

// IPAllocationId is a struct representing the Resource ID for a I P Allocation
type IPAllocationId struct {
	SubscriptionId    string
	ResourceGroupName string
	IpAllocationName  string
}

// NewIPAllocationID returns a new IPAllocationId struct
func NewIPAllocationID(subscriptionId string, resourceGroupName string, ipAllocationName string) IPAllocationId {
	return IPAllocationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		IpAllocationName:  ipAllocationName,
	}
}

// ParseIPAllocationID parses 'input' into a IPAllocationId
func ParseIPAllocationID(input string) (*IPAllocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(IPAllocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IPAllocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IpAllocationName, ok = parsed.Parsed["ipAllocationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ipAllocationName", *parsed)
	}

	return &id, nil
}

// ParseIPAllocationIDInsensitively parses 'input' case-insensitively into a IPAllocationId
// note: this method should only be used for API response data and not user input
func ParseIPAllocationIDInsensitively(input string) (*IPAllocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(IPAllocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IPAllocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.IpAllocationName, ok = parsed.Parsed["ipAllocationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "ipAllocationName", *parsed)
	}

	return &id, nil
}

// ValidateIPAllocationID checks that 'input' can be parsed as a I P Allocation ID
func ValidateIPAllocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIPAllocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted I P Allocation ID
func (id IPAllocationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/ipAllocations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IpAllocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this I P Allocation ID
func (id IPAllocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticIpAllocations", "ipAllocations", "ipAllocations"),
		resourceids.UserSpecifiedSegment("ipAllocationName", "ipAllocationValue"),
	}
}

// String returns a human-readable description of this I P Allocation ID
func (id IPAllocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ip Allocation Name: %q", id.IpAllocationName),
	}
	return fmt.Sprintf("I P Allocation (%s)", strings.Join(components, "\n"))
}
