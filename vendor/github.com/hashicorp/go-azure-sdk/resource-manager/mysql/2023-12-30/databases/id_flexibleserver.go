package databases

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FlexibleServerId{})
}

var _ resourceids.ResourceId = &FlexibleServerId{}

// FlexibleServerId is a struct representing the Resource ID for a Flexible Server
type FlexibleServerId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
}

// NewFlexibleServerID returns a new FlexibleServerId struct
func NewFlexibleServerID(subscriptionId string, resourceGroupName string, flexibleServerName string) FlexibleServerId {
	return FlexibleServerId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
	}
}

// ParseFlexibleServerID parses 'input' into a FlexibleServerId
func ParseFlexibleServerID(input string) (*FlexibleServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FlexibleServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FlexibleServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFlexibleServerIDInsensitively parses 'input' case-insensitively into a FlexibleServerId
// note: this method should only be used for API response data and not user input
func ParseFlexibleServerIDInsensitively(input string) (*FlexibleServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FlexibleServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FlexibleServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FlexibleServerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	return nil
}

// ValidateFlexibleServerID checks that 'input' can be parsed as a Flexible Server ID
func ValidateFlexibleServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFlexibleServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Flexible Server ID
func (id FlexibleServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Flexible Server ID
func (id FlexibleServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMySQL", "Microsoft.DBforMySQL", "Microsoft.DBforMySQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
	}
}

// String returns a human-readable description of this Flexible Server ID
func (id FlexibleServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
	}
	return fmt.Sprintf("Flexible Server (%s)", strings.Join(components, "\n"))
}
