package actiongroupsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ActionGroupId{}

// ActionGroupId is a struct representing the Resource ID for a Action Group
type ActionGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	ActionGroupName   string
}

// NewActionGroupID returns a new ActionGroupId struct
func NewActionGroupID(subscriptionId string, resourceGroupName string, actionGroupName string) ActionGroupId {
	return ActionGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ActionGroupName:   actionGroupName,
	}
}

// ParseActionGroupID parses 'input' into a ActionGroupId
func ParseActionGroupID(input string) (*ActionGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ActionGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ActionGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ActionGroupName, ok = parsed.Parsed["actionGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "actionGroupName", *parsed)
	}

	return &id, nil
}

// ParseActionGroupIDInsensitively parses 'input' case-insensitively into a ActionGroupId
// note: this method should only be used for API response data and not user input
func ParseActionGroupIDInsensitively(input string) (*ActionGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ActionGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ActionGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ActionGroupName, ok = parsed.Parsed["actionGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "actionGroupName", *parsed)
	}

	return &id, nil
}

// ValidateActionGroupID checks that 'input' can be parsed as a Action Group ID
func ValidateActionGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseActionGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Action Group ID
func (id ActionGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/actionGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ActionGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Action Group ID
func (id ActionGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticActionGroups", "actionGroups", "actionGroups"),
		resourceids.UserSpecifiedSegment("actionGroupName", "actionGroupValue"),
	}
}

// String returns a human-readable description of this Action Group ID
func (id ActionGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Action Group Name: %q", id.ActionGroupName),
	}
	return fmt.Sprintf("Action Group (%s)", strings.Join(components, "\n"))
}
