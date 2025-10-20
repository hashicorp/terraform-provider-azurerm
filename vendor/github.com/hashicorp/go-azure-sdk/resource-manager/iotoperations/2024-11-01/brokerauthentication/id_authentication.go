package brokerauthentication

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AuthenticationId{})
}

var _ resourceids.ResourceId = &AuthenticationId{}

// AuthenticationId is a struct representing the Resource ID for a Authentication
type AuthenticationId struct {
	SubscriptionId     string
	ResourceGroupName  string
	InstanceName       string
	BrokerName         string
	AuthenticationName string
}

// NewAuthenticationID returns a new AuthenticationId struct
func NewAuthenticationID(subscriptionId string, resourceGroupName string, instanceName string, brokerName string, authenticationName string) AuthenticationId {
	return AuthenticationId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		InstanceName:       instanceName,
		BrokerName:         brokerName,
		AuthenticationName: authenticationName,
	}
}

// ParseAuthenticationID parses 'input' into a AuthenticationId
func ParseAuthenticationID(input string) (*AuthenticationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthenticationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthenticationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthenticationIDInsensitively parses 'input' case-insensitively into a AuthenticationId
// note: this method should only be used for API response data and not user input
func ParseAuthenticationIDInsensitively(input string) (*AuthenticationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthenticationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthenticationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthenticationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.InstanceName, ok = input.Parsed["instanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceName", input)
	}

	if id.BrokerName, ok = input.Parsed["brokerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "brokerName", input)
	}

	if id.AuthenticationName, ok = input.Parsed["authenticationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authenticationName", input)
	}

	return nil
}

// ValidateAuthenticationID checks that 'input' can be parsed as a Authentication ID
func ValidateAuthenticationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthenticationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Authentication ID
func (id AuthenticationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s/brokers/%s/authentications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.InstanceName, id.BrokerName, id.AuthenticationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Authentication ID
func (id AuthenticationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftIoTOperations", "Microsoft.IoTOperations", "Microsoft.IoTOperations"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceName"),
		resourceids.StaticSegment("staticBrokers", "brokers", "brokers"),
		resourceids.UserSpecifiedSegment("brokerName", "brokerName"),
		resourceids.StaticSegment("staticAuthentications", "authentications", "authentications"),
		resourceids.UserSpecifiedSegment("authenticationName", "authenticationName"),
	}
}

// String returns a human-readable description of this Authentication ID
func (id AuthenticationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
		fmt.Sprintf("Broker Name: %q", id.BrokerName),
		fmt.Sprintf("Authentication Name: %q", id.AuthenticationName),
	}
	return fmt.Sprintf("Authentication (%s)", strings.Join(components, "\n"))
}
