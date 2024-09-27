package scalingplan

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScalingPlanId{})
}

var _ resourceids.ResourceId = &ScalingPlanId{}

// ScalingPlanId is a struct representing the Resource ID for a Scaling Plan
type ScalingPlanId struct {
	SubscriptionId    string
	ResourceGroupName string
	ScalingPlanName   string
}

// NewScalingPlanID returns a new ScalingPlanId struct
func NewScalingPlanID(subscriptionId string, resourceGroupName string, scalingPlanName string) ScalingPlanId {
	return ScalingPlanId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ScalingPlanName:   scalingPlanName,
	}
}

// ParseScalingPlanID parses 'input' into a ScalingPlanId
func ParseScalingPlanID(input string) (*ScalingPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScalingPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScalingPlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScalingPlanIDInsensitively parses 'input' case-insensitively into a ScalingPlanId
// note: this method should only be used for API response data and not user input
func ParseScalingPlanIDInsensitively(input string) (*ScalingPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScalingPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScalingPlanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScalingPlanId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ScalingPlanName, ok = input.Parsed["scalingPlanName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scalingPlanName", input)
	}

	return nil
}

// ValidateScalingPlanID checks that 'input' can be parsed as a Scaling Plan ID
func ValidateScalingPlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScalingPlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scaling Plan ID
func (id ScalingPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/scalingPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ScalingPlanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scaling Plan ID
func (id ScalingPlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticScalingPlans", "scalingPlans", "scalingPlans"),
		resourceids.UserSpecifiedSegment("scalingPlanName", "scalingPlanName"),
	}
}

// String returns a human-readable description of this Scaling Plan ID
func (id ScalingPlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Scaling Plan Name: %q", id.ScalingPlanName),
	}
	return fmt.Sprintf("Scaling Plan (%s)", strings.Join(components, "\n"))
}
