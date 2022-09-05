package scalingplan

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScalingPlanId{}

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
	parser := resourceids.NewParserFromResourceIdType(ScalingPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScalingPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ScalingPlanName, ok = parsed.Parsed["scalingPlanName"]; !ok {
		return nil, fmt.Errorf("the segment 'scalingPlanName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScalingPlanIDInsensitively parses 'input' case-insensitively into a ScalingPlanId
// note: this method should only be used for API response data and not user input
func ParseScalingPlanIDInsensitively(input string) (*ScalingPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScalingPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScalingPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ScalingPlanName, ok = parsed.Parsed["scalingPlanName"]; !ok {
		return nil, fmt.Errorf("the segment 'scalingPlanName' was not found in the resource id %q", input)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("scalingPlanName", "scalingPlanValue"),
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
