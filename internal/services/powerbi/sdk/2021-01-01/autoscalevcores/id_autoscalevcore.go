package autoscalevcores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AutoScaleVCoreId{}

// AutoScaleVCoreId is a struct representing the Resource ID for a Auto Scale V Core
type AutoScaleVCoreId struct {
	SubscriptionId    string
	ResourceGroupName string
	VcoreName         string
}

// NewAutoScaleVCoreID returns a new AutoScaleVCoreId struct
func NewAutoScaleVCoreID(subscriptionId string, resourceGroupName string, vcoreName string) AutoScaleVCoreId {
	return AutoScaleVCoreId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VcoreName:         vcoreName,
	}
}

// ParseAutoScaleVCoreID parses 'input' into a AutoScaleVCoreId
func ParseAutoScaleVCoreID(input string) (*AutoScaleVCoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(AutoScaleVCoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AutoScaleVCoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VcoreName, ok = parsed.Parsed["vcoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'vcoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAutoScaleVCoreIDInsensitively parses 'input' case-insensitively into a AutoScaleVCoreId
// note: this method should only be used for API response data and not user input
func ParseAutoScaleVCoreIDInsensitively(input string) (*AutoScaleVCoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(AutoScaleVCoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AutoScaleVCoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VcoreName, ok = parsed.Parsed["vcoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'vcoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAutoScaleVCoreID checks that 'input' can be parsed as a Auto Scale V Core ID
func ValidateAutoScaleVCoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutoScaleVCoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Auto Scale V Core ID
func (id AutoScaleVCoreId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PowerBIDedicated/autoScaleVCores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VcoreName)
}

// Segments returns a slice of Resource ID Segments which comprise this Auto Scale V Core ID
func (id AutoScaleVCoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftPowerBIDedicated", "Microsoft.PowerBIDedicated", "Microsoft.PowerBIDedicated"),
		resourceids.StaticSegment("autoScaleVCores", "autoScaleVCores", "autoScaleVCores"),
		resourceids.UserSpecifiedSegment("vcoreName", "vcoreValue"),
	}
}

// String returns a human-readable description of this Auto Scale V Core ID
func (id AutoScaleVCoreId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vcore Name: %q", id.VcoreName),
	}
	return fmt.Sprintf("Auto Scale V Core (%s)", strings.Join(components, "\n"))
}
