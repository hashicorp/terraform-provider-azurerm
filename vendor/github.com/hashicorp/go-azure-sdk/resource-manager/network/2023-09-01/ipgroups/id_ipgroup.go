package ipgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &IPGroupId{}

// IPGroupId is a struct representing the Resource ID for a I P Group
type IPGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	IpGroupName       string
}

// NewIPGroupID returns a new IPGroupId struct
func NewIPGroupID(subscriptionId string, resourceGroupName string, ipGroupName string) IPGroupId {
	return IPGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		IpGroupName:       ipGroupName,
	}
}

// ParseIPGroupID parses 'input' into a IPGroupId
func ParseIPGroupID(input string) (*IPGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IPGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IPGroupId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIPGroupIDInsensitively parses 'input' case-insensitively into a IPGroupId
// note: this method should only be used for API response data and not user input
func ParseIPGroupIDInsensitively(input string) (*IPGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IPGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IPGroupId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IPGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.IpGroupName, ok = input.Parsed["ipGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ipGroupName", input)
	}

	return nil
}

// ValidateIPGroupID checks that 'input' can be parsed as a I P Group ID
func ValidateIPGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIPGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted I P Group ID
func (id IPGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/ipGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IpGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this I P Group ID
func (id IPGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticIpGroups", "ipGroups", "ipGroups"),
		resourceids.UserSpecifiedSegment("ipGroupName", "ipGroupValue"),
	}
}

// String returns a human-readable description of this I P Group ID
func (id IPGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ip Group Name: %q", id.IpGroupName),
	}
	return fmt.Sprintf("I P Group (%s)", strings.Join(components, "\n"))
}
