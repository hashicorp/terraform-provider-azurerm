package applicationgatewayprivateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationGatewayPrivateEndpointConnectionId{})
}

var _ resourceids.ResourceId = &ApplicationGatewayPrivateEndpointConnectionId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationGatewayPrivateEndpointConnectionIDInsensitively parses 'input' case-insensitively into a ApplicationGatewayPrivateEndpointConnectionId
// note: this method should only be used for API response data and not user input
func ParseApplicationGatewayPrivateEndpointConnectionIDInsensitively(input string) (*ApplicationGatewayPrivateEndpointConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayPrivateEndpointConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayPrivateEndpointConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationGatewayPrivateEndpointConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ApplicationGatewayName, ok = input.Parsed["applicationGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationGatewayName", input)
	}

	if id.PrivateEndpointConnectionName, ok = input.Parsed["privateEndpointConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointConnectionName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("applicationGatewayName", "applicationGatewayName"),
		resourceids.StaticSegment("staticPrivateEndpointConnections", "privateEndpointConnections", "privateEndpointConnections"),
		resourceids.UserSpecifiedSegment("privateEndpointConnectionName", "privateEndpointConnectionName"),
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
