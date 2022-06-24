package operationstatus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ProviderOperationStatuId{}

// ProviderOperationStatuId is a struct representing the Resource ID for a Provider Operation Statu
type ProviderOperationStatuId struct {
	SubscriptionId    string
	ResourceGroupName string
	OperationId       string
}

// NewProviderOperationStatuID returns a new ProviderOperationStatuId struct
func NewProviderOperationStatuID(subscriptionId string, resourceGroupName string, operationId string) ProviderOperationStatuId {
	return ProviderOperationStatuId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		OperationId:       operationId,
	}
}

// ParseProviderOperationStatuID parses 'input' into a ProviderOperationStatuId
func ParseProviderOperationStatuID(input string) (*ProviderOperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderOperationStatuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderOperationStatuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProviderOperationStatuIDInsensitively parses 'input' case-insensitively into a ProviderOperationStatuId
// note: this method should only be used for API response data and not user input
func ParseProviderOperationStatuIDInsensitively(input string) (*ProviderOperationStatuId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderOperationStatuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderOperationStatuId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProviderOperationStatuID checks that 'input' can be parsed as a Provider Operation Statu ID
func ValidateProviderOperationStatuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderOperationStatuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Operation Statu ID
func (id ProviderOperationStatuId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/operationStatus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Operation Statu ID
func (id ProviderOperationStatuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticOperationStatus", "operationStatus", "operationStatus"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Provider Operation Statu ID
func (id ProviderOperationStatuId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Provider Operation Statu (%s)", strings.Join(components, "\n"))
}
