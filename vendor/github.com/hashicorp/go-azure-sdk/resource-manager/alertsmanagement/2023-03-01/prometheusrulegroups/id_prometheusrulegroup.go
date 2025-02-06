package prometheusrulegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrometheusRuleGroupId{})
}

var _ resourceids.ResourceId = &PrometheusRuleGroupId{}

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
	parser := resourceids.NewParserFromResourceIdType(&PrometheusRuleGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrometheusRuleGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrometheusRuleGroupIDInsensitively parses 'input' case-insensitively into a PrometheusRuleGroupId
// note: this method should only be used for API response data and not user input
func ParsePrometheusRuleGroupIDInsensitively(input string) (*PrometheusRuleGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrometheusRuleGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrometheusRuleGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrometheusRuleGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrometheusRuleGroupName, ok = input.Parsed["prometheusRuleGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "prometheusRuleGroupName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("prometheusRuleGroupName", "prometheusRuleGroupName"),
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
