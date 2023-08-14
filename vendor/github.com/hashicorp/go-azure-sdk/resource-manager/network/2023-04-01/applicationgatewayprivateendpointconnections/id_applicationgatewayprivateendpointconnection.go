package applicationgatewayprivateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApplicationGatewayPrivateEndpointConnectionId{}

// ApplicationGatewayPrivateEndpointConnectionId is a struct representing the Resource ID for a Application Gateway Private Endpoint Connection
type ApplicationGatewayPrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroupName             string
	ApplicationGatewayName        string
	PrivateEndpointConnectionName string
}

// NewApplicationGatewayPrivateEndpointConnectionID returns a new ApplicationGatewayPrivateEndpointConnectionId struct
func NewApplicationGatewayPrivateEndpointConnectionID(subscriptionId string, resourceGroupName string, applicationGatewayName string, privateEndpointConnectionName string) ApplicationGatewayPrivateEndpointConnectionId {
	return ApplicationGatewayPrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		ApplicationGatewayName:        applicationGatewayName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

// ParseApplicationGatewayPrivateEndpointConnectionID parses 'input' into a ApplicationGatewayPrivateEndpointConnectionId
func ParseApplicationGatewayPrivateEndpointConnectionID(input string) (*ApplicationGatewayPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationGatewayPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationGatewayPrivateEndpointConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ApplicationGatewayName, ok = parsed.Parsed["applicationGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applicationGatewayName", *parsed)
	}

	if id.PrivateEndpointConnectionName, ok = parsed.Parsed["privateEndpointConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", *parsed)
	}

	return &id, nil
}

// ParseApplicationGatewayPrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a ApplicationGatewayPrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParseApplicationGatewayPrivateEndpointConnectionIDInsensitively(input string) (*ApplicationGatewayPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationGatewayPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationGatewayPrivateEndpointConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ApplicationGatewayName, ok = parsed.Parsed["applicationGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applicationGatewayName", *parsed)
	}

	if id.PrivateEndpointConnectionName, ok = parsed.Parsed["privateEndpointConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateApplicationGatewayPrivateEndpointConnectionID checks that 'input' can be parsed as a Application Gateway Private Endpoint Connection ID
func ValidateApplicationGatewayPrivateEndpointConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationGatewayPrivateEndpointConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Gateway Private Endpoint Connection ID
func (id ApplicationGatewayPrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGatewayName, id.PrivateEndpointConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Gateway Private Endpoint Connection ID
func (id ApplicationGatewayPrivateEndpointConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationGateways", "applicationGateways", "applicationGateways"),
		resourceids.UserSpecifiedSegment("applicationGatewayName", "applicationGatewayValue"),
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionValue"),
	}
}

// String returns a human-readable description of this Application Gateway Private Endpoint Connection ID
func (id ApplicationGatewayPrivateEndpointConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Gateway Name: %q", id.ApplicationGatewayName),
		fmt.Sprintf("Private Endpoint Connection Name: %q", id.PrivateEndpointConnectionName),
	}
	return fmt.Sprintf("Application Gateway Private Endpoint Connection (%s)", strings.Join(components, "\n"))
}
