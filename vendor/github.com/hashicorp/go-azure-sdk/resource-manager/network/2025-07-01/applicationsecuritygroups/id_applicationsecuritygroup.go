package applicationsecuritygroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationSecurityGroupId{})
}

var _ resourceids.ResourceId = &ApplicationSecurityGroupId{}

// ApplicationSecurityGroupId is a struct representing the Resource ID for a Application Security Group
type ApplicationSecurityGroupId struct {
	SubscriptionId               string
	ResourceGroupName            string
	ApplicationSecurityGroupName string
}

// NewApplicationSecurityGroupID returns a new ApplicationSecurityGroupId struct
func NewApplicationSecurityGroupID(subscriptionId string, resourceGroupName string, applicationSecurityGroupName string) ApplicationSecurityGroupId {
	return ApplicationSecurityGroupId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		ApplicationSecurityGroupName: applicationSecurityGroupName,
	}
}

// ParseApplicationSecurityGroupID parses 'input' into a ApplicationSecurityGroupId
func ParseApplicationSecurityGroupID(input string) (*ApplicationSecurityGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationSecurityGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationSecurityGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationSecurityGroupIDInsensitively parses 'input' case-insensitively into a ApplicationSecurityGroupId
// note: this method should only be used for API response data and not user input
func ParseApplicationSecurityGroupIDInsensitively(input string) (*ApplicationSecurityGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationSecurityGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationSecurityGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationSecurityGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ApplicationSecurityGroupName, ok = input.Parsed["applicationSecurityGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationSecurityGroupName", input)
	}

	return nil
}

// ValidateApplicationSecurityGroupID checks that 'input' can be parsed as a Application Security Group ID
func ValidateApplicationSecurityGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationSecurityGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Security Group ID
func (id ApplicationSecurityGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationSecurityGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationSecurityGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Security Group ID
func (id ApplicationSecurityGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationSecurityGroups", "applicationSecurityGroups", "applicationSecurityGroups"),
		resourceids.UserSpecifiedSegment("applicationSecurityGroupName", "applicationSecurityGroupName"),
	}
}

// String returns a human-readable description of this Application Security Group ID
func (id ApplicationSecurityGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Security Group Name: %q", id.ApplicationSecurityGroupName),
	}
	return fmt.Sprintf("Application Security Group (%s)", strings.Join(components, "\n"))
}
