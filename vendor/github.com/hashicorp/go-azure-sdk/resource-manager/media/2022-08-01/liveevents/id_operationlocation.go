package liveevents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OperationLocationId{}

// OperationLocationId is a struct representing the Resource ID for a Operation Location
type OperationLocationId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	LiveEventName     string
	OperationId       string
}

// NewOperationLocationID returns a new OperationLocationId struct
func NewOperationLocationID(subscriptionId string, resourceGroupName string, accountName string, liveEventName string, operationId string) OperationLocationId {
	return OperationLocationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		LiveEventName:     liveEventName,
		OperationId:       operationId,
	}
}

// ParseOperationLocationID parses 'input' into a OperationLocationId
func ParseOperationLocationID(input string) (*OperationLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationLocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationLocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.LiveEventName, ok = parsed.Parsed["liveEventName"]; !ok {
		return nil, fmt.Errorf("the segment 'liveEventName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseOperationLocationIDInsensitively parses 'input' case-insensitively into a OperationLocationId
// note: this method should only be used for API response data and not user input
func ParseOperationLocationIDInsensitively(input string) (*OperationLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(OperationLocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OperationLocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.LiveEventName, ok = parsed.Parsed["liveEventName"]; !ok {
		return nil, fmt.Errorf("the segment 'liveEventName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateOperationLocationID checks that 'input' can be parsed as a Operation Location ID
func ValidateOperationLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOperationLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Operation Location ID
func (id OperationLocationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/liveEvents/%s/operationLocations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.LiveEventName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Operation Location ID
func (id OperationLocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticLiveEvents", "liveEvents", "liveEvents"),
		resourceids.UserSpecifiedSegment("liveEventName", "liveEventValue"),
		resourceids.StaticSegment("staticOperationLocations", "operationLocations", "operationLocations"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Operation Location ID
func (id OperationLocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Live Event Name: %q", id.LiveEventName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Operation Location (%s)", strings.Join(components, "\n"))
}
