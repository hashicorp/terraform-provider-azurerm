package metricalerts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = MetricAlertId{}

// MetricAlertId is a struct representing the Resource ID for a Metric Alert
type MetricAlertId struct {
	SubscriptionId    string
	ResourceGroupName string
	MetricAlertName   string
}

// NewMetricAlertID returns a new MetricAlertId struct
func NewMetricAlertID(subscriptionId string, resourceGroupName string, metricAlertName string) MetricAlertId {
	return MetricAlertId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MetricAlertName:   metricAlertName,
	}
}

// ParseMetricAlertID parses 'input' into a MetricAlertId
func ParseMetricAlertID(input string) (*MetricAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(MetricAlertId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MetricAlertId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MetricAlertName, ok = parsed.Parsed["metricAlertName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "metricAlertName", *parsed)
	}

	return &id, nil
}

// ParseMetricAlertIDInsensitively parses 'input' case-insensitively into a MetricAlertId
// note: this method should only be used for API response data and not user input
func ParseMetricAlertIDInsensitively(input string) (*MetricAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(MetricAlertId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MetricAlertId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MetricAlertName, ok = parsed.Parsed["metricAlertName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "metricAlertName", *parsed)
	}

	return &id, nil
}

// ValidateMetricAlertID checks that 'input' can be parsed as a Metric Alert ID
func ValidateMetricAlertID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMetricAlertID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Metric Alert ID
func (id MetricAlertId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/metricAlerts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MetricAlertName)
}

// Segments returns a slice of Resource ID Segments which comprise this Metric Alert ID
func (id MetricAlertId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticMetricAlerts", "metricAlerts", "metricAlerts"),
		resourceids.UserSpecifiedSegment("metricAlertName", "metricAlertValue"),
	}
}

// String returns a human-readable description of this Metric Alert ID
func (id MetricAlertId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Metric Alert Name: %q", id.MetricAlertName),
	}
	return fmt.Sprintf("Metric Alert (%s)", strings.Join(components, "\n"))
}
