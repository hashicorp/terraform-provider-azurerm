package apigatewayconfigconnection

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConfigConnectionId{})
}

var _ resourceids.ResourceId = &ConfigConnectionId{}

// ConfigConnectionId is a struct representing the Resource ID for a Config Connection
type ConfigConnectionId struct {
	SubscriptionId       string
	ResourceGroupName    string
	GatewayName          string
	ConfigConnectionName string
}

// NewConfigConnectionID returns a new ConfigConnectionId struct
func NewConfigConnectionID(subscriptionId string, resourceGroupName string, gatewayName string, configConnectionName string) ConfigConnectionId {
	return ConfigConnectionId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		GatewayName:          gatewayName,
		ConfigConnectionName: configConnectionName,
	}
}

// ParseConfigConnectionID parses 'input' into a ConfigConnectionId
func ParseConfigConnectionID(input string) (*ConfigConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigConnectionIDInsensitively parses 'input' case-insensitively into a ConfigConnectionId
// note: this method should only be used for API response data and not user input
func ParseConfigConnectionIDInsensitively(input string) (*ConfigConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GatewayName, ok = input.Parsed["gatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gatewayName", input)
	}

	if id.ConfigConnectionName, ok = input.Parsed["configConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configConnectionName", input)
	}

	return nil
}

// ValidateConfigConnectionID checks that 'input' can be parsed as a Config Connection ID
func ValidateConfigConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Config Connection ID
func (id ConfigConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/gateways/%s/configConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GatewayName, id.ConfigConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Config Connection ID
func (id ConfigConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayName"),
		resourceids.StaticSegment("staticConfigConnections", "configConnections", "configConnections"),
		resourceids.UserSpecifiedSegment("configConnectionName", "configConnectionName"),
	}
}

// String returns a human-readable description of this Config Connection ID
func (id ConfigConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Gateway Name: %q", id.GatewayName),
		fmt.Sprintf("Config Connection Name: %q", id.ConfigConnectionName),
	}
	return fmt.Sprintf("Config Connection (%s)", strings.Join(components, "\n"))
}
