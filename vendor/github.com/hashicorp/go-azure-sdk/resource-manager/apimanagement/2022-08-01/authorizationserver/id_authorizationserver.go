package authorizationserver

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AuthorizationServerId{})
}

var _ resourceids.ResourceId = &AuthorizationServerId{}

// AuthorizationServerId is a struct representing the Resource ID for a Authorization Server
type AuthorizationServerId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ServiceName             string
	AuthorizationServerName string
}

// NewAuthorizationServerID returns a new AuthorizationServerId struct
func NewAuthorizationServerID(subscriptionId string, resourceGroupName string, serviceName string, authorizationServerName string) AuthorizationServerId {
	return AuthorizationServerId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ServiceName:             serviceName,
		AuthorizationServerName: authorizationServerName,
	}
}

// ParseAuthorizationServerID parses 'input' into a AuthorizationServerId
func ParseAuthorizationServerID(input string) (*AuthorizationServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthorizationServerIDInsensitively parses 'input' case-insensitively into a AuthorizationServerId
// note: this method should only be used for API response data and not user input
func ParseAuthorizationServerIDInsensitively(input string) (*AuthorizationServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthorizationServerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AuthorizationServerName, ok = input.Parsed["authorizationServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationServerName", input)
	}

	return nil
}

// ValidateAuthorizationServerID checks that 'input' can be parsed as a Authorization Server ID
func ValidateAuthorizationServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthorizationServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Authorization Server ID
func (id AuthorizationServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/authorizationServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.AuthorizationServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Authorization Server ID
func (id AuthorizationServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticAuthorizationServers", "authorizationServers", "authorizationServers"),
		resourceids.UserSpecifiedSegment("authorizationServerName", "authorizationServerName"),
	}
}

// String returns a human-readable description of this Authorization Server ID
func (id AuthorizationServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Authorization Server Name: %q", id.AuthorizationServerName),
	}
	return fmt.Sprintf("Authorization Server (%s)", strings.Join(components, "\n"))
}
