// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DedicatedHostGroupId{}

// DedicatedHostGroupId is a struct representing the Resource ID for a Dedicated Host Group
type DedicatedHostGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostGroupName     string
}

// NewDedicatedHostGroupID returns a new HostGroupId struct
func NewDedicatedHostGroupID(subscriptionId string, resourceGroupName string, hostGroupName string) DedicatedHostGroupId {
	return DedicatedHostGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostGroupName:     hostGroupName,
	}
}

// ParseDedicatedHostGroupID parses 'input' into a DedicatedHostGroupId
func ParseDedicatedHostGroupID(input string) (*DedicatedHostGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(DedicatedHostGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DedicatedHostGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HostGroupName, ok = parsed.Parsed["hostGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hostGroupName", *parsed)
	}

	return &id, nil
}

// ParseDedicatedHostGroupIDInsensitively parses 'input' case-insensitively into a DedicatedHostGroupId
// note: this method should only be used for API response data and not user input
func ParseDedicatedHostGroupIDInsensitively(input string) (*DedicatedHostGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(DedicatedHostGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DedicatedHostGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.HostGroupName, ok = parsed.Parsed["hostGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "hostGroupName", *parsed)
	}

	return &id, nil
}

// ValidateDedicatedHostGroupID checks that 'input' can be parsed as a Dedicated Host Group ID
func ValidateDedicatedHostGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDedicatedHostGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dedicated Host Group ID
func (id DedicatedHostGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dedicated Host Group ID
func (id DedicatedHostGroupId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Dedicated Host Group ID
func (id DedicatedHostGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Group Name: %q", id.HostGroupName),
	}
	return fmt.Sprintf("Dedicated Host Group (%s)", strings.Join(components, "\n"))
}
