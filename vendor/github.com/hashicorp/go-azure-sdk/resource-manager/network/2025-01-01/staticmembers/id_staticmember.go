package staticmembers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StaticMemberId{})
}

var _ resourceids.ResourceId = &StaticMemberId{}

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
	parser := resourceids.NewParserFromResourceIdType(&StaticMemberId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticMemberId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStaticMemberIDInsensitively parses 'input' case-insensitively into a StaticMemberId
// note: this method should only be used for API response data and not user input
func ParseStaticMemberIDInsensitively(input string) (*StaticMemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StaticMemberId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StaticMemberId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StaticMemberId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkManagerName, ok = input.Parsed["networkManagerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", input)
	}

	if id.NetworkGroupName, ok = input.Parsed["networkGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkGroupName", input)
	}

	if id.StaticMemberName, ok = input.Parsed["staticMemberName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticMemberName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
		resourceids.StaticSegment("staticNetworkGroups", "networkGroups", "networkGroups"),
		resourceids.UserSpecifiedSegment("networkGroupName", "networkGroupName"),
		resourceids.StaticSegment("staticStaticMembers", "staticMembers", "staticMembers"),
		resourceids.UserSpecifiedSegment("staticMemberName", "staticMemberName"),
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
