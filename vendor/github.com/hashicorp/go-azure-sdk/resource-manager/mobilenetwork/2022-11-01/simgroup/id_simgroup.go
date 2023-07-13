package simgroup

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SimGroupId{}

// SimGroupId is a struct representing the Resource ID for a Sim Group
type SimGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	SimGroupName      string
}

// NewSimGroupID returns a new SimGroupId struct
func NewSimGroupID(subscriptionId string, resourceGroupName string, simGroupName string) SimGroupId {
	return SimGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SimGroupName:      simGroupName,
	}
}

// ParseSimGroupID parses 'input' into a SimGroupId
func ParseSimGroupID(input string) (*SimGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SimGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SimGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SimGroupName, ok = parsed.Parsed["simGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "simGroupName", *parsed)
	}

	return &id, nil
}

// ParseSimGroupIDInsensitively parses 'input' case-insensitively into a SimGroupId
// note: this method should only be used for API response data and not user input
func ParseSimGroupIDInsensitively(input string) (*SimGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SimGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SimGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SimGroupName, ok = parsed.Parsed["simGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "simGroupName", *parsed)
	}

	return &id, nil
}

// ValidateSimGroupID checks that 'input' can be parsed as a Sim Group ID
func ValidateSimGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSimGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sim Group ID
func (id SimGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/simGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SimGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sim Group ID
func (id SimGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticSimGroups", "simGroups", "simGroups"),
		resourceids.UserSpecifiedSegment("simGroupName", "simGroupValue"),
	}
}

// String returns a human-readable description of this Sim Group ID
func (id SimGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sim Group Name: %q", id.SimGroupName),
	}
	return fmt.Sprintf("Sim Group (%s)", strings.Join(components, "\n"))
}
