package metricalerts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MetricAlertId{})
}

var _ resourceids.ResourceId = &MetricAlertId{}

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
	parser := resourceids.NewParserFromResourceIdType(&MetricAlertId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MetricAlertId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMetricAlertIDInsensitively parses 'input' case-insensitively into a MetricAlertId
// note: this method should only be used for API response data and not user input
func ParseMetricAlertIDInsensitively(input string) (*MetricAlertId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MetricAlertId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MetricAlertId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MetricAlertId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MetricAlertName, ok = input.Parsed["metricAlertName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "metricAlertName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("metricAlertName", "metricAlertName"),
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
