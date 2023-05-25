package prometheusrulegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrometheusRuleGroupId{}

// PrometheusRuleGroupId is a struct representing the Resource ID for a Prometheus Rule Group
type PrometheusRuleGroupId struct {
	SubscriptionId          string
	ResourceGroupName       string
	PrometheusRuleGroupName string
}

// NewPrometheusRuleGroupID returns a new PrometheusRuleGroupId struct
func NewPrometheusRuleGroupID(subscriptionId string, resourceGroupName string, prometheusRuleGroupName string) PrometheusRuleGroupId {
	return PrometheusRuleGroupId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		PrometheusRuleGroupName: prometheusRuleGroupName,
	}
}

// ParsePrometheusRuleGroupID parses 'input' into a PrometheusRuleGroupId
func ParsePrometheusRuleGroupID(input string) (*PrometheusRuleGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrometheusRuleGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrometheusRuleGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrometheusRuleGroupName, ok = parsed.Parsed["prometheusRuleGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "prometheusRuleGroupName", *parsed)
	}

	return &id, nil
}

// ParsePrometheusRuleGroupIDInsensitively parses 'input' case-insensitively into a PrometheusRuleGroupId
// note: this method should only be used for API response data and not user input
func ParsePrometheusRuleGroupIDInsensitively(input string) (*PrometheusRuleGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrometheusRuleGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrometheusRuleGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrometheusRuleGroupName, ok = parsed.Parsed["prometheusRuleGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "prometheusRuleGroupName", *parsed)
	}

	return &id, nil
}

// ValidatePrometheusRuleGroupID checks that 'input' can be parsed as a Prometheus Rule Group ID
func ValidatePrometheusRuleGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrometheusRuleGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Prometheus Rule Group ID
func (id PrometheusRuleGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AlertsManagement/prometheusRuleGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrometheusRuleGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Prometheus Rule Group ID
func (id PrometheusRuleGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAlertsManagement", "Microsoft.AlertsManagement", "Microsoft.AlertsManagement"),
		resourceids.StaticSegment("staticPrometheusRuleGroups", "prometheusRuleGroups", "prometheusRuleGroups"),
		resourceids.UserSpecifiedSegment("prometheusRuleGroupName", "prometheusRuleGroupValue"),
	}
}

// String returns a human-readable description of this Prometheus Rule Group ID
func (id PrometheusRuleGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Prometheus Rule Group Name: %q", id.PrometheusRuleGroupName),
	}
	return fmt.Sprintf("Prometheus Rule Group (%s)", strings.Join(components, "\n"))
}
