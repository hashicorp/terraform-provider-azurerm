package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkTraceId{})
}

var _ resourceids.ResourceId = &NetworkTraceId{}

// NetworkTraceId is a struct representing the Resource ID for a Network Trace
type NetworkTraceId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	OperationId       string
}

// NewNetworkTraceID returns a new NetworkTraceId struct
func NewNetworkTraceID(subscriptionId string, resourceGroupName string, siteName string, operationId string) NetworkTraceId {
	return NetworkTraceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		OperationId:       operationId,
	}
}

// ParseNetworkTraceID parses 'input' into a NetworkTraceId
func ParseNetworkTraceID(input string) (*NetworkTraceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkTraceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkTraceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkTraceIDInsensitively parses 'input' case-insensitively into a NetworkTraceId
// note: this method should only be used for API response data and not user input
func ParseNetworkTraceIDInsensitively(input string) (*NetworkTraceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkTraceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkTraceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkTraceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.OperationId, ok = input.Parsed["operationId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "operationId", input)
	}

	return nil
}

// ValidateNetworkTraceID checks that 'input' can be parsed as a Network Trace ID
func ValidateNetworkTraceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkTraceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Trace ID
func (id NetworkTraceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/networkTrace/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.OperationId)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Trace ID
func (id NetworkTraceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticNetworkTrace", "networkTrace", "networkTrace"),
		resourceids.UserSpecifiedSegment("operationId", "operationId"),
	}
}

// String returns a human-readable description of this Network Trace ID
func (id NetworkTraceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Operation: %q", id.OperationId),
	}
	return fmt.Sprintf("Network Trace (%s)", strings.Join(components, "\n"))
}
