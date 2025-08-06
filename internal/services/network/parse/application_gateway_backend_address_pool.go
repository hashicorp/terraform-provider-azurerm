// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ApplicationGatewayBackendAddressPoolId{}

// ApplicationGatewayBackendAddressPoolId is a struct representing the Resource ID for a Compilation Job
type ApplicationGatewayBackendAddressPoolId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ApplicationGatewayName string
	BackendAddressPoolName string
}

// NewApplicationGatewayBackendAddressPoolID returns a new ApplicationGatewayBackendAddressPoolId struct
func NewApplicationGatewayBackendAddressPoolID(subscriptionId string, resourceGroupName string, applicationGatewayName string, backendAddressPoolName string) ApplicationGatewayBackendAddressPoolId {
	return ApplicationGatewayBackendAddressPoolId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ApplicationGatewayName: applicationGatewayName,
		BackendAddressPoolName: backendAddressPoolName,
	}
}

// ApplicationGatewayBackendAddressPoolID parses 'input' into a ApplicationGatewayBackendAddressPoolId
func ApplicationGatewayBackendAddressPoolID(input string) (*ApplicationGatewayBackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayBackendAddressPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayBackendAddressPoolId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationGatewayBackendAddressPoolIDInsensitively parses 'input' case-insensitively into a ApplicationGatewayBackendAddressPoolId
// note: this method should only be used for API response data and not user input
func ParseApplicationGatewayBackendAddressPoolIDInsensitively(input string) (*ApplicationGatewayBackendAddressPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationGatewayBackendAddressPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationGatewayBackendAddressPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationGatewayBackendAddressPoolId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BackendAddressPoolName, ok = input.Parsed["backendAddressPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backendAddressPoolName", input)
	}

	return nil
}

// ValidateApplicationGatewayBackendAddressPoolID checks that 'input' can be parsed as a Compilation Job ID
func ValidateApplicationGatewayBackendAddressPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ApplicationGatewayBackendAddressPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Compilation Job ID
func (id ApplicationGatewayBackendAddressPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/backendAddressPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGatewayName, id.BackendAddressPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Services I P Configuration ID
func (id ApplicationGatewayBackendAddressPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticApplicationGateways", "applicationGateways", "applicationGateways"),
		resourceids.UserSpecifiedSegment("applicationGatewayName", "applicationGatewaysValue"),
		resourceids.StaticSegment("staticBackendAddressPools", "backendAddressPools", "backendAddressPools"),
		resourceids.UserSpecifiedSegment("backendAddressPoolName", "backendAddressPoolNameValue"),
	}
}

// String returns a human-readable description of this Backend Address Pool ID
func (id ApplicationGatewayBackendAddressPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Gateway Name: %q", id.ApplicationGatewayName),
		fmt.Sprintf("Backend Address Pool: %q", id.BackendAddressPoolName),
	}
	return fmt.Sprintf("Application Gateway Backend Address Pool (%s)", strings.Join(components, "\n"))
}
