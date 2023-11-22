package staticmembers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = StaticMemberId{}

// StaticMemberId is a struct representing the Resource ID for a Static Member
type StaticMemberId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkManagerName string
	NetworkGroupName   string
	StaticMemberName   string
}

// NewStaticMemberID returns a new StaticMemberId struct
func NewStaticMemberID(subscriptionId string, resourceGroupName string, networkManagerName string, networkGroupName string, staticMemberName string) StaticMemberId {
	return StaticMemberId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkManagerName: networkManagerName,
		NetworkGroupName:   networkGroupName,
		StaticMemberName:   staticMemberName,
	}
}

// ParseStaticMemberID parses 'input' into a StaticMemberId
func ParseStaticMemberID(input string) (*StaticMemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(StaticMemberId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StaticMemberId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.NetworkGroupName, ok = parsed.Parsed["networkGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkGroupName", *parsed)
	}

	if id.StaticMemberName, ok = parsed.Parsed["staticMemberName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "staticMemberName", *parsed)
	}

	return &id, nil
}

// ParseStaticMemberIDInsensitively parses 'input' case-insensitively into a StaticMemberId
// note: this method should only be used for API response data and not user input
func ParseStaticMemberIDInsensitively(input string) (*StaticMemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(StaticMemberId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StaticMemberId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.NetworkGroupName, ok = parsed.Parsed["networkGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkGroupName", *parsed)
	}

	if id.StaticMemberName, ok = parsed.Parsed["staticMemberName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "staticMemberName", *parsed)
	}

	return &id, nil
}

// ValidateStaticMemberID checks that 'input' can be parsed as a Static Member ID
func ValidateStaticMemberID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStaticMemberID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Static Member ID
func (id StaticMemberId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/networkGroups/%s/staticMembers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
}

// Segments returns a slice of Resource ID Segments which comprise this Static Member ID
func (id StaticMemberId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerValue"),
		resourceids.StaticSegment("staticNetworkGroups", "networkGroups", "networkGroups"),
		resourceids.UserSpecifiedSegment("networkGroupName", "networkGroupValue"),
		resourceids.StaticSegment("staticStaticMembers", "staticMembers", "staticMembers"),
		resourceids.UserSpecifiedSegment("staticMemberName", "staticMemberValue"),
	}
}

// String returns a human-readable description of this Static Member ID
func (id StaticMemberId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Network Group Name: %q", id.NetworkGroupName),
		fmt.Sprintf("Static Member Name: %q", id.StaticMemberName),
	}
	return fmt.Sprintf("Static Member (%s)", strings.Join(components, "\n"))
}
