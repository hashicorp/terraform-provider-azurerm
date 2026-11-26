package expressroutelagmembers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MemberId{})
}

var _ resourceids.ResourceId = &MemberId{}

// MemberId is a struct representing the Resource ID for a Member
type MemberId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ExpressRouteLagName string
	LinkName            string
	MemberName          string
}

// NewMemberID returns a new MemberId struct
func NewMemberID(subscriptionId string, resourceGroupName string, expressRouteLagName string, linkName string, memberName string) MemberId {
	return MemberId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ExpressRouteLagName: expressRouteLagName,
		LinkName:            linkName,
		MemberName:          memberName,
	}
}

// ParseMemberID parses 'input' into a MemberId
func ParseMemberID(input string) (*MemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MemberId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MemberId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMemberIDInsensitively parses 'input' case-insensitively into a MemberId
// note: this method should only be used for API response data and not user input
func ParseMemberIDInsensitively(input string) (*MemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MemberId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MemberId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MemberId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteLagName, ok = input.Parsed["expressRouteLagName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteLagName", input)
	}

	if id.LinkName, ok = input.Parsed["linkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkName", input)
	}

	if id.MemberName, ok = input.Parsed["memberName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "memberName", input)
	}

	return nil
}

// ValidateMemberID checks that 'input' can be parsed as a Member ID
func ValidateMemberID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMemberID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Member ID
func (id MemberId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteLags/%s/links/%s/members/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteLagName, id.LinkName, id.MemberName)
}

// Segments returns a slice of Resource ID Segments which comprise this Member ID
func (id MemberId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteLags", "expressRouteLags", "expressRouteLags"),
		resourceids.UserSpecifiedSegment("expressRouteLagName", "expressRouteLagName"),
		resourceids.StaticSegment("staticLinks", "links", "links"),
		resourceids.UserSpecifiedSegment("linkName", "linkName"),
		resourceids.StaticSegment("staticMembers", "members", "members"),
		resourceids.UserSpecifiedSegment("memberName", "memberName"),
	}
}

// String returns a human-readable description of this Member ID
func (id MemberId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Lag Name: %q", id.ExpressRouteLagName),
		fmt.Sprintf("Link Name: %q", id.LinkName),
		fmt.Sprintf("Member Name: %q", id.MemberName),
	}
	return fmt.Sprintf("Member (%s)", strings.Join(components, "\n"))
}
