package grafanaresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = GrafanaId{}

// GrafanaId is a struct representing the Resource ID for a Grafana
type GrafanaId struct {
	SubscriptionId    string
	ResourceGroupName string
	GrafanaName       string
}

// NewGrafanaID returns a new GrafanaId struct
func NewGrafanaID(subscriptionId string, resourceGroupName string, grafanaName string) GrafanaId {
	return GrafanaId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GrafanaName:       grafanaName,
	}
}

// ParseGrafanaID parses 'input' into a GrafanaId
func ParseGrafanaID(input string) (*GrafanaId, error) {
	parser := resourceids.NewParserFromResourceIdType(GrafanaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GrafanaId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GrafanaName, ok = parsed.Parsed["grafanaName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "grafanaName", *parsed)
	}

	return &id, nil
}

// ParseGrafanaIDInsensitively parses 'input' case-insensitively into a GrafanaId
// note: this method should only be used for API response data and not user input
func ParseGrafanaIDInsensitively(input string) (*GrafanaId, error) {
	parser := resourceids.NewParserFromResourceIdType(GrafanaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GrafanaId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.GrafanaName, ok = parsed.Parsed["grafanaName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "grafanaName", *parsed)
	}

	return &id, nil
}

// ValidateGrafanaID checks that 'input' can be parsed as a Grafana ID
func ValidateGrafanaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGrafanaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Grafana ID
func (id GrafanaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Dashboard/grafana/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GrafanaName)
}

// Segments returns a slice of Resource ID Segments which comprise this Grafana ID
func (id GrafanaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDashboard", "Microsoft.Dashboard", "Microsoft.Dashboard"),
		resourceids.StaticSegment("staticGrafana", "grafana", "grafana"),
		resourceids.UserSpecifiedSegment("grafanaName", "grafanaValue"),
	}
}

// String returns a human-readable description of this Grafana ID
func (id GrafanaId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Grafana Name: %q", id.GrafanaName),
	}
	return fmt.Sprintf("Grafana (%s)", strings.Join(components, "\n"))
}
