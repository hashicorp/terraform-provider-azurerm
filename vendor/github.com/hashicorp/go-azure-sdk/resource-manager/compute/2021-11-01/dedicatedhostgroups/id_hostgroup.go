package dedicatedhostgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = HostGroupId{}

// HostGroupId is a struct representing the Resource ID for a Host Group
type HostGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostGroupName     string
}

// NewHostGroupID returns a new HostGroupId struct
func NewHostGroupID(subscriptionId string, resourceGroupName string, hostGroupName string) HostGroupId {
	return HostGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostGroupName:     hostGroupName,
	}
}

// ParseHostGroupID parses 'input' into a HostGroupId
func ParseHostGroupID(input string) (*HostGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostGroupId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHostGroupIDInsensitively parses 'input' case-insensitively into a HostGroupId
// note: this method should only be used for API response data and not user input
func ParseHostGroupIDInsensitively(input string) (*HostGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostGroupId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HostGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostGroupName, ok = input.Parsed["hostGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostGroupName", input)
	}

	return nil
}

// ValidateHostGroupID checks that 'input' can be parsed as a Host Group ID
func ValidateHostGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Host Group ID
func (id HostGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Host Group ID
func (id HostGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticHostGroups", "hostGroups", "hostGroups"),
		resourceids.UserSpecifiedSegment("hostGroupName", "hostGroupValue"),
	}
}

// String returns a human-readable description of this Host Group ID
func (id HostGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Group Name: %q", id.HostGroupName),
	}
	return fmt.Sprintf("Host Group (%s)", strings.Join(components, "\n"))
}
