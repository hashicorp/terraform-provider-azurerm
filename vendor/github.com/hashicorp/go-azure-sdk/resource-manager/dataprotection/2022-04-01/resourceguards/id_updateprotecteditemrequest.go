package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = UpdateProtectedItemRequestId{}

// UpdateProtectedItemRequestId is a struct representing the Resource ID for a Update Protected Item Request
type UpdateProtectedItemRequestId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ResourceGuardsName string
	RequestName        string
}

// NewUpdateProtectedItemRequestID returns a new UpdateProtectedItemRequestId struct
func NewUpdateProtectedItemRequestID(subscriptionId string, resourceGroupName string, resourceGuardsName string, requestName string) UpdateProtectedItemRequestId {
	return UpdateProtectedItemRequestId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ResourceGuardsName: resourceGuardsName,
		RequestName:        requestName,
	}
}

// ParseUpdateProtectedItemRequestID parses 'input' into a UpdateProtectedItemRequestId
func ParseUpdateProtectedItemRequestID(input string) (*UpdateProtectedItemRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateProtectedItemRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateProtectedItemRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceGuardsName, ok = parsed.Parsed["resourceGuardsName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGuardsName' was not found in the resource id %q", input)
	}

	if id.RequestName, ok = parsed.Parsed["requestName"]; !ok {
		return nil, fmt.Errorf("the segment 'requestName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseUpdateProtectedItemRequestIDInsensitively parses 'input' case-insensitively into a UpdateProtectedItemRequestId
// note: this method should only be used for API response data and not user input
func ParseUpdateProtectedItemRequestIDInsensitively(input string) (*UpdateProtectedItemRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateProtectedItemRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateProtectedItemRequestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceGuardsName, ok = parsed.Parsed["resourceGuardsName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGuardsName' was not found in the resource id %q", input)
	}

	if id.RequestName, ok = parsed.Parsed["requestName"]; !ok {
		return nil, fmt.Errorf("the segment 'requestName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateUpdateProtectedItemRequestID checks that 'input' can be parsed as a Update Protected Item Request ID
func ValidateUpdateProtectedItemRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateProtectedItemRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update Protected Item Request ID
func (id UpdateProtectedItemRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/updateProtectedItemRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardsName, id.RequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update Protected Item Request ID
func (id UpdateProtectedItemRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardsName", "resourceGuardsValue"),
		resourceids.StaticSegment("staticUpdateProtectedItemRequests", "updateProtectedItemRequests", "updateProtectedItemRequests"),
		resourceids.UserSpecifiedSegment("requestName", "requestValue"),
	}
}

// String returns a human-readable description of this Update Protected Item Request ID
func (id UpdateProtectedItemRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guards Name: %q", id.ResourceGuardsName),
		fmt.Sprintf("Request Name: %q", id.RequestName),
	}
	return fmt.Sprintf("Update Protected Item Request (%s)", strings.Join(components, "\n"))
}
