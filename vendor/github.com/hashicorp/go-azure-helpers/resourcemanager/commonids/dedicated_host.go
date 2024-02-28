// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &DedicatedHostId{}

// DedicatedHostId is a struct representing the Resource ID for a Dedicated Host
type DedicatedHostId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostGroupName     string
	HostName          string
}

// NewDedicatedHostID returns a new DedicatedHostId struct
func NewDedicatedHostID(subscriptionId string, resourceGroupName string, hostGroupName string, hostName string) DedicatedHostId {
	return DedicatedHostId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostGroupName:     hostGroupName,
		HostName:          hostName,
	}
}

// ParseDedicatedHostID parses 'input' into a DedicatedHostId
func ParseDedicatedHostID(input string) (*DedicatedHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DedicatedHostId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DedicatedHostId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDedicatedHostIDInsensitively parses 'input' case-insensitively into a DedicatedHostId
// note: this method should only be used for API response data and not user input
func ParseDedicatedHostIDInsensitively(input string) (*DedicatedHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DedicatedHostId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DedicatedHostId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ValidateDedicatedHostID checks that 'input' can be parsed as a Dedicated Host ID
func ValidateDedicatedHostID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDedicatedHostID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

func (id *DedicatedHostId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.HostName, ok = input.Parsed["hostName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostName", input)
	}

	return nil
}

// ID returns the formatted Dedicated Host ID
func (id DedicatedHostId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s/hosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostGroupName, id.HostName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dedicated Host ID
func (id DedicatedHostId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticHostGroups", "hostGroups", "hostGroups"),
		resourceids.UserSpecifiedSegment("hostGroupName", "hostGroupValue"),
		resourceids.StaticSegment("staticHosts", "hosts", "hosts"),
		resourceids.UserSpecifiedSegment("hostName", "hostValue"),
	}
}

// String returns a human-readable description of this DedicatedHost ID
func (id DedicatedHostId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Group Name: %q", id.HostGroupName),
		fmt.Sprintf("Host Name: %q", id.HostName),
	}
	return fmt.Sprintf("Dedicated Host (%s)", strings.Join(components, "\n"))
}
