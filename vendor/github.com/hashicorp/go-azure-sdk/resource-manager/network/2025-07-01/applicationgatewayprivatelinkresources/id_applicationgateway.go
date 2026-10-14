package applicationgatewayprivatelinkresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationGatewayId{})
}

var _ resourceids.ResourceId = &ApplicationGatewayId{}

// ApplicationGatewayId is a struct representing the Resource ID for a Application Gateway
type ApplicationGatewayId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ApplicationGatewayName string
}

// NewApplicationGatewayID returns a new ApplicationGatewayId struct
func NewApplicationGatewayID(subscriptionId string, resourceGroupName string, applicationGatewayName string) ApplicationGatewayId {
	return ApplicationGatewayId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ApplicationGatewayName: applicationGatewayName,
	}
}

// ParseApplicationGatewayID parses 'input' into a ApplicationGatewayId
func ParseApplicationGatewayID(input string) (*ApplicationGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationGatewayIDInsensitively parses 'input' case-insensitively into a ApplicationGatewayId
// note: this method should only be used for API response data and not user input
func ParseApplicationGatewayIDInsensitively(input string) (*ApplicationGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationGatewayId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateApplicationGatewayID checks that 'input' can be parsed as a Application Gateway ID
func ValidateApplicationGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Gateway ID
func (id ApplicationGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Gateway ID
func (id ApplicationGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationGateways", "applicationGateways", "applicationGateways"),
		resourceids.UserSpecifiedSegment("applicationGatewayName", "applicationGatewayName"),
	}
}

// String returns a human-readable description of this Application Gateway ID
func (id ApplicationGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Gateway Name: %q", id.ApplicationGatewayName),
	}
	return fmt.Sprintf("Application Gateway (%s)", strings.Join(components, "\n"))
}
