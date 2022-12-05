package encodings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = TransformId{}

// TransformId is a struct representing the Resource ID for a Transform
type TransformId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	TransformName     string
}

// NewTransformID returns a new TransformId struct
func NewTransformID(subscriptionId string, resourceGroupName string, accountName string, transformName string) TransformId {
	return TransformId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		TransformName:     transformName,
	}
}

// ParseTransformID parses 'input' into a TransformId
func ParseTransformID(input string) (*TransformId, error) {
	parser := resourceids.NewParserFromResourceIdType(TransformId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TransformId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.TransformName, ok = parsed.Parsed["transformName"]; !ok {
		return nil, fmt.Errorf("the segment 'transformName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseTransformIDInsensitively parses 'input' case-insensitively into a TransformId
// note: this method should only be used for API response data and not user input
func ParseTransformIDInsensitively(input string) (*TransformId, error) {
	parser := resourceids.NewParserFromResourceIdType(TransformId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TransformId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.TransformName, ok = parsed.Parsed["transformName"]; !ok {
		return nil, fmt.Errorf("the segment 'transformName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateTransformID checks that 'input' can be parsed as a Transform ID
func ValidateTransformID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTransformID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Transform ID
func (id TransformId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/transforms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.TransformName)
}

// Segments returns a slice of Resource ID Segments which comprise this Transform ID
func (id TransformId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticTransforms", "transforms", "transforms"),
		resourceids.UserSpecifiedSegment("transformName", "transformValue"),
	}
}

// String returns a human-readable description of this Transform ID
func (id TransformId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Transform Name: %q", id.TransformName),
	}
	return fmt.Sprintf("Transform (%s)", strings.Join(components, "\n"))
}
