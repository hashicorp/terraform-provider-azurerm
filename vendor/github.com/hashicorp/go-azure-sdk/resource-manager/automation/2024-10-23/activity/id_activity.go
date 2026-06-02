package activity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ActivityId{})
}

var _ resourceids.ResourceId = &ActivityId{}

// ActivityId is a struct representing the Resource ID for a Activity
type ActivityId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	ModuleName            string
	ActivityName          string
}

// NewActivityID returns a new ActivityId struct
func NewActivityID(subscriptionId string, resourceGroupName string, automationAccountName string, moduleName string, activityName string) ActivityId {
	return ActivityId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		ModuleName:            moduleName,
		ActivityName:          activityName,
	}
}

// ParseActivityID parses 'input' into a ActivityId
func ParseActivityID(input string) (*ActivityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ActivityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ActivityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseActivityIDInsensitively parses 'input' case-insensitively into a ActivityId
// note: this method should only be used for API response data and not user input
func ParseActivityIDInsensitively(input string) (*ActivityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ActivityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ActivityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ActivityId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutomationAccountName, ok = input.Parsed["automationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", input)
	}

	if id.ModuleName, ok = input.Parsed["moduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "moduleName", input)
	}

	if id.ActivityName, ok = input.Parsed["activityName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "activityName", input)
	}

	return nil
}

// ValidateActivityID checks that 'input' can be parsed as a Activity ID
func ValidateActivityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseActivityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Activity ID
func (id ActivityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/modules/%s/activities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.ModuleName, id.ActivityName)
}

// Segments returns a slice of Resource ID Segments which comprise this Activity ID
func (id ActivityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticModules", "modules", "modules"),
		resourceids.UserSpecifiedSegment("moduleName", "moduleName"),
		resourceids.StaticSegment("staticActivities", "activities", "activities"),
		resourceids.UserSpecifiedSegment("activityName", "activityName"),
	}
}

// String returns a human-readable description of this Activity ID
func (id ActivityId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Module Name: %q", id.ModuleName),
		fmt.Sprintf("Activity Name: %q", id.ActivityName),
	}
	return fmt.Sprintf("Activity (%s)", strings.Join(components, "\n"))
}
