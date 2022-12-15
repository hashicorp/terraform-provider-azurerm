package streamingendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = StreamingEndpointOperationId{}

// StreamingEndpointOperationId is a struct representing the Resource ID for a Streaming Endpoint Operation
type StreamingEndpointOperationId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	OperationId       string
}

// NewStreamingEndpointOperationID returns a new StreamingEndpointOperationId struct
func NewStreamingEndpointOperationID(subscriptionId string, resourceGroupName string, accountName string, operationId string) StreamingEndpointOperationId {
	return StreamingEndpointOperationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		OperationId:       operationId,
	}
}

// ParseStreamingEndpointOperationID parses 'input' into a StreamingEndpointOperationId
func ParseStreamingEndpointOperationID(input string) (*StreamingEndpointOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingEndpointOperationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingEndpointOperationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseStreamingEndpointOperationIDInsensitively parses 'input' case-insensitively into a StreamingEndpointOperationId
// note: this method should only be used for API response data and not user input
func ParseStreamingEndpointOperationIDInsensitively(input string) (*StreamingEndpointOperationId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingEndpointOperationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingEndpointOperationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.OperationId, ok = parsed.Parsed["operationId"]; !ok {
		return nil, fmt.Errorf("the segment 'operationId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateStreamingEndpointOperationID checks that 'input' can be parsed as a Streaming Endpoint Operation ID
func ValidateStreamingEndpointOperationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStreamingEndpointOperationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Streaming Endpoint Operation ID
func (id StreamingEndpointOperationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/streamingEndpointOperations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Streaming Endpoint Operation ID
func (id StreamingEndpointOperationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticStreamingEndpointOperations", "streamingEndpointOperations", "streamingEndpointOperations"),
		resourceids.UserSpecifiedSegment("operationId", "operationIdValue"),
	}
}

// String returns a human-readable description of this Streaming Endpoint Operation ID
func (id StreamingEndpointOperationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Streaming Endpoint Operation (%s)", strings.Join(components, "\n"))
}
