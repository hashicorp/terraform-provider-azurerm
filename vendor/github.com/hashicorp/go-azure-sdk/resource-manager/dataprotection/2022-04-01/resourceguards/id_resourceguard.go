package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ResourceGuardId{}

// ResourceGuardId is a struct representing the Resource ID for a Resource Guard
type ResourceGuardId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ResourceGuardsName string
}

// NewResourceGuardID returns a new ResourceGuardId struct
func NewResourceGuardID(subscriptionId string, resourceGroupName string, resourceGuardsName string) ResourceGuardId {
	return ResourceGuardId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ResourceGuardsName: resourceGuardsName,
	}
}

// ParseResourceGuardID parses 'input' into a ResourceGuardId
func ParseResourceGuardID(input string) (*ResourceGuardId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceGuardId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceGuardId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceGuardsName, ok = parsed.Parsed["resourceGuardsName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGuardsName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseResourceGuardIDInsensitively parses 'input' case-insensitively into a ResourceGuardId
// note: this method should only be used for API response data and not user input
func ParseResourceGuardIDInsensitively(input string) (*ResourceGuardId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceGuardId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceGuardId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceGuardsName, ok = parsed.Parsed["resourceGuardsName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGuardsName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateResourceGuardID checks that 'input' can be parsed as a Resource Guard ID
func ValidateResourceGuardID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceGuardID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Guard ID
func (id ResourceGuardId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardsName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Guard ID
func (id ResourceGuardId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardsName", "resourceGuardsValue"),
	}
}

// String returns a human-readable description of this Resource Guard ID
func (id ResourceGuardId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guards Name: %q", id.ResourceGuardsName),
	}
	return fmt.Sprintf("Resource Guard (%s)", strings.Join(components, "\n"))
}
