package activitylogalertsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ActivityLogAlertId{}

// ActivityLogAlertId is a struct representing the Resource ID for a Activity Log Alert
type ActivityLogAlertId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ActivityLogAlertName string
}

// NewActivityLogAlertID returns a new ActivityLogAlertId struct
func NewActivityLogAlertID(subscriptionId string, resourceGroupName string, activityLogAlertName string) ActivityLogAlertId {
	return ActivityLogAlertId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ActivityLogAlertName: activityLogAlertName,
	}
}

// ParseActivityLogAlertID parses 'input' into a ActivityLogAlertId
func ParseActivityLogAlertID(input string) (*ActivityLogAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(ActivityLogAlertId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ActivityLogAlertId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ActivityLogAlertName, ok = parsed.Parsed["activityLogAlertName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "activityLogAlertName", *parsed)
	}

	return &id, nil
}

// ParseActivityLogAlertIDInsensitively parses 'input' case-insensitively into a ActivityLogAlertId
// note: this method should only be used for API response data and not user input
func ParseActivityLogAlertIDInsensitively(input string) (*ActivityLogAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(ActivityLogAlertId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ActivityLogAlertId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ActivityLogAlertName, ok = parsed.Parsed["activityLogAlertName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "activityLogAlertName", *parsed)
	}

	return &id, nil
}

// ValidateActivityLogAlertID checks that 'input' can be parsed as a Activity Log Alert ID
func ValidateActivityLogAlertID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseActivityLogAlertID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Activity Log Alert ID
func (id ActivityLogAlertId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/activityLogAlerts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ActivityLogAlertName)
}

// Segments returns a slice of Resource ID Segments which comprise this Activity Log Alert ID
func (id ActivityLogAlertId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticActivityLogAlerts", "activityLogAlerts", "activityLogAlerts"),
		resourceids.UserSpecifiedSegment("activityLogAlertName", "activityLogAlertValue"),
	}
}

// String returns a human-readable description of this Activity Log Alert ID
func (id ActivityLogAlertId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Activity Log Alert Name: %q", id.ActivityLogAlertName),
	}
	return fmt.Sprintf("Activity Log Alert (%s)", strings.Join(components, "\n"))
}
