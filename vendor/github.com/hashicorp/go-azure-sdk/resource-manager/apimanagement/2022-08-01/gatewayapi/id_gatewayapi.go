package gatewayapi

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GatewayApiId{})
}

var _ resourceids.ResourceId = &GatewayApiId{}

// GatewayApiId is a struct representing the Resource ID for a Gateway Api
type GatewayApiId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	GatewayId         string
	ApiId             string
}

// NewGatewayApiID returns a new GatewayApiId struct
func NewGatewayApiID(subscriptionId string, resourceGroupName string, serviceName string, gatewayId string, apiId string) GatewayApiId {
	return GatewayApiId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		GatewayId:         gatewayId,
		ApiId:             apiId,
	}
}

// ParseGatewayApiID parses 'input' into a GatewayApiId
func ParseGatewayApiID(input string) (*GatewayApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GatewayApiId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GatewayApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGatewayApiIDInsensitively parses 'input' case-insensitively into a GatewayApiId
// note: this method should only be used for API response data and not user input
func ParseGatewayApiIDInsensitively(input string) (*GatewayApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GatewayApiId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GatewayApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GatewayApiId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.GatewayId, ok = input.Parsed["gatewayId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gatewayId", input)
	}

	if id.ApiId, ok = input.Parsed["apiId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiId", input)
	}

	return nil
}

// ValidateGatewayApiID checks that 'input' can be parsed as a Gateway Api ID
func ValidateGatewayApiID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGatewayApiID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gateway Api ID
func (id GatewayApiId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/apis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId, id.ApiId)
}

// Segments returns a slice of Resource ID Segments which comprise this Gateway Api ID
func (id GatewayApiId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayId", "gatewayId"),
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiId", "apiId"),
	}
}

// String returns a human-readable description of this Gateway Api ID
func (id GatewayApiId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Gateway: %q", id.GatewayId),
		fmt.Sprintf("Api: %q", id.ApiId),
	}
	return fmt.Sprintf("Gateway Api (%s)", strings.Join(components, "\n"))
}
