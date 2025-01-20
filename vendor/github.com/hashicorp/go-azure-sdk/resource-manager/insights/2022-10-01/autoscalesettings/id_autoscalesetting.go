package autoscalesettings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutoScaleSettingId{})
}

var _ resourceids.ResourceId = &AutoScaleSettingId{}

// AutoScaleSettingId is a struct representing the Resource ID for a Auto Scale Setting
type AutoScaleSettingId struct {
	SubscriptionId       string
	ResourceGroupName    string
	AutoScaleSettingName string
}

// NewAutoScaleSettingID returns a new AutoScaleSettingId struct
func NewAutoScaleSettingID(subscriptionId string, resourceGroupName string, autoScaleSettingName string) AutoScaleSettingId {
	return AutoScaleSettingId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		AutoScaleSettingName: autoScaleSettingName,
	}
}

// ParseAutoScaleSettingID parses 'input' into a AutoScaleSettingId
func ParseAutoScaleSettingID(input string) (*AutoScaleSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoScaleSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoScaleSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutoScaleSettingIDInsensitively parses 'input' case-insensitively into a AutoScaleSettingId
// note: this method should only be used for API response data and not user input
func ParseAutoScaleSettingIDInsensitively(input string) (*AutoScaleSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutoScaleSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutoScaleSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutoScaleSettingId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutoScaleSettingName, ok = input.Parsed["autoScaleSettingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autoScaleSettingName", input)
	}

	return nil
}

// ValidateAutoScaleSettingID checks that 'input' can be parsed as a Auto Scale Setting ID
func ValidateAutoScaleSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutoScaleSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Auto Scale Setting ID
func (id AutoScaleSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/autoScaleSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutoScaleSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Auto Scale Setting ID
func (id AutoScaleSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticAutoScaleSettings", "autoScaleSettings", "autoScaleSettings"),
		resourceids.UserSpecifiedSegment("autoScaleSettingName", "autoScaleSettingName"),
	}
}

// String returns a human-readable description of this Auto Scale Setting ID
func (id AutoScaleSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Auto Scale Setting Name: %q", id.AutoScaleSettingName),
	}
	return fmt.Sprintf("Auto Scale Setting (%s)", strings.Join(components, "\n"))
}
