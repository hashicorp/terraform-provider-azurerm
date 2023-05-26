package alertruletemplates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AlertRuleTemplateId{}

// AlertRuleTemplateId is a struct representing the Resource ID for a Alert Rule Template
type AlertRuleTemplateId struct {
	SubscriptionId      string
	ResourceGroupName   string
	WorkspaceName       string
	AlertRuleTemplateId string
}

// NewAlertRuleTemplateID returns a new AlertRuleTemplateId struct
func NewAlertRuleTemplateID(subscriptionId string, resourceGroupName string, workspaceName string, alertRuleTemplateId string) AlertRuleTemplateId {
	return AlertRuleTemplateId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		WorkspaceName:       workspaceName,
		AlertRuleTemplateId: alertRuleTemplateId,
	}
}

// ParseAlertRuleTemplateID parses 'input' into a AlertRuleTemplateId
func ParseAlertRuleTemplateID(input string) (*AlertRuleTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(AlertRuleTemplateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AlertRuleTemplateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.AlertRuleTemplateId, ok = parsed.Parsed["alertRuleTemplateId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "alertRuleTemplateId", *parsed)
	}

	return &id, nil
}

// ParseAlertRuleTemplateIDInsensitively parses 'input' case-insensitively into a AlertRuleTemplateId
// note: this method should only be used for API response data and not user input
func ParseAlertRuleTemplateIDInsensitively(input string) (*AlertRuleTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(AlertRuleTemplateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AlertRuleTemplateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.AlertRuleTemplateId, ok = parsed.Parsed["alertRuleTemplateId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "alertRuleTemplateId", *parsed)
	}

	return &id, nil
}

// ValidateAlertRuleTemplateID checks that 'input' can be parsed as a Alert Rule Template ID
func ValidateAlertRuleTemplateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAlertRuleTemplateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Alert Rule Template ID
func (id AlertRuleTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/alertRuleTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.AlertRuleTemplateId)
}

// Segments returns a slice of Resource ID Segments which comprise this Alert Rule Template ID
func (id AlertRuleTemplateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticAlertRuleTemplates", "alertRuleTemplates", "alertRuleTemplates"),
		resourceids.UserSpecifiedSegment("alertRuleTemplateId", "alertRuleTemplateIdValue"),
	}
}

// String returns a human-readable description of this Alert Rule Template ID
func (id AlertRuleTemplateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Alert Rule Template: %q", id.AlertRuleTemplateId),
	}
	return fmt.Sprintf("Alert Rule Template (%s)", strings.Join(components, "\n"))
}
