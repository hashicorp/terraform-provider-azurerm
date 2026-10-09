package msiximage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&HostPoolId{})
}

var _ resourceids.ResourceId = &HostPoolId{}

// HostPoolId is a struct representing the Resource ID for a Host Pool
type HostPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostPoolName      string
}

// NewHostPoolID returns a new HostPoolId struct
func NewHostPoolID(subscriptionId string, resourceGroupName string, hostPoolName string) HostPoolId {
	return HostPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostPoolName:      hostPoolName,
	}
}

// ParseHostPoolID parses 'input' into a HostPoolId
func ParseHostPoolID(input string) (*HostPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HostPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHostPoolIDInsensitively parses 'input' case-insensitively into a HostPoolId
// note: this method should only be used for API response data and not user input
func ParseHostPoolIDInsensitively(input string) (*HostPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HostPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HostPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HostPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostPoolName, ok = input.Parsed["hostPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostPoolName", input)
	}

	return nil
}

// ValidateHostPoolID checks that 'input' can be parsed as a Host Pool ID
func ValidateHostPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Host Pool ID
func (id HostPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Host Pool ID
func (id HostPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticHostPools", "hostPools", "hostPools"),
		resourceids.UserSpecifiedSegment("hostPoolName", "hostPoolName"),
	}
}

// String returns a human-readable description of this Host Pool ID
func (id HostPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Pool Name: %q", id.HostPoolName),
	}
	return fmt.Sprintf("Host Pool (%s)", strings.Join(components, "\n"))
}
