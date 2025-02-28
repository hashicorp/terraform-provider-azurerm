package gateway

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ServiceGatewayId{})
}

var _ resourceids.ResourceId = &ServiceGatewayId{}

// ServiceGatewayId is a struct representing the Resource ID for a Service Gateway
type ServiceGatewayId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	GatewayId         string
}

// NewServiceGatewayID returns a new ServiceGatewayId struct
func NewServiceGatewayID(subscriptionId string, resourceGroupName string, serviceName string, gatewayId string) ServiceGatewayId {
	return ServiceGatewayId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		GatewayId:         gatewayId,
	}
}

// ParseServiceGatewayID parses 'input' into a ServiceGatewayId
func ParseServiceGatewayID(input string) (*ServiceGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ServiceGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ServiceGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseServiceGatewayIDInsensitively parses 'input' case-insensitively into a ServiceGatewayId
// note: this method should only be used for API response data and not user input
func ParseServiceGatewayIDInsensitively(input string) (*ServiceGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ServiceGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ServiceGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ServiceGatewayId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateServiceGatewayID checks that 'input' can be parsed as a Service Gateway ID
func ValidateServiceGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseServiceGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Service Gateway ID
func (id ServiceGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId)
}

// Segments returns a slice of Resource ID Segments which comprise this Service Gateway ID
func (id ServiceGatewayId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Service Gateway ID
func (id ServiceGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Gateway: %q", id.GatewayId),
	}
	return fmt.Sprintf("Service Gateway (%s)", strings.Join(components, "\n"))
}
