package grafanaresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GrafanaId{})
}

var _ resourceids.ResourceId = &GrafanaId{}

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
	parser := resourceids.NewParserFromResourceIdType(&GrafanaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GrafanaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGrafanaIDInsensitively parses 'input' case-insensitively into a GrafanaId
// note: this method should only be used for API response data and not user input
func ParseGrafanaIDInsensitively(input string) (*GrafanaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GrafanaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GrafanaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GrafanaId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GrafanaName, ok = input.Parsed["grafanaName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "grafanaName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("grafanaName", "grafanaName"),
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
