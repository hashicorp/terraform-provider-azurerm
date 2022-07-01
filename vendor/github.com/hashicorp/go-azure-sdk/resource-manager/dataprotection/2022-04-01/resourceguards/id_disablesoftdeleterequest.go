package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DisableSoftDeleteRequestId{}

// DisableSoftDeleteRequestId is a struct representing the Resource ID for a Disable Soft Delete Request
type DisableSoftDeleteRequestId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ResourceGuardsName string
	RequestName        string
}

// NewDisableSoftDeleteRequestID returns a new DisableSoftDeleteRequestId struct
func NewDisableSoftDeleteRequestID(subscriptionId string, resourceGroupName string, resourceGuardsName string, requestName string) DisableSoftDeleteRequestId {
	return DisableSoftDeleteRequestId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ResourceGuardsName: resourceGuardsName,
		RequestName:        requestName,
	}
}

// ParseDisableSoftDeleteRequestID parses 'input' into a DisableSoftDeleteRequestId
func ParseDisableSoftDeleteRequestID(input string) (*DisableSoftDeleteRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(DisableSoftDeleteRequestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DisableSoftDeleteRequestId{}

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

// ParseDisableSoftDeleteRequestIDInsensitively parses 'input' case-insensitively into a DisableSoftDeleteRequestId
// note: this method should only be used for API response data and not user input
func ParseDisableSoftDeleteRequestIDInsensitively(input string) (*DisableSoftDeleteRequestId, error) {
	parser := resourceids.NewParserFromResourceIdType(DisableSoftDeleteRequestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DisableSoftDeleteRequestId{}

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

// ValidateDisableSoftDeleteRequestID checks that 'input' can be parsed as a Disable Soft Delete Request ID
func ValidateDisableSoftDeleteRequestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDisableSoftDeleteRequestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Disable Soft Delete Request ID
func (id DisableSoftDeleteRequestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s/disableSoftDeleteRequests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardsName, id.RequestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Disable Soft Delete Request ID
func (id DisableSoftDeleteRequestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardsName", "resourceGuardsValue"),
		resourceids.StaticSegment("staticDisableSoftDeleteRequests", "disableSoftDeleteRequests", "disableSoftDeleteRequests"),
		resourceids.UserSpecifiedSegment("requestName", "requestValue"),
	}
}

// String returns a human-readable description of this Disable Soft Delete Request ID
func (id DisableSoftDeleteRequestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guards Name: %q", id.ResourceGuardsName),
		fmt.Sprintf("Request Name: %q", id.RequestName),
	}
	return fmt.Sprintf("Disable Soft Delete Request (%s)", strings.Join(components, "\n"))
}
