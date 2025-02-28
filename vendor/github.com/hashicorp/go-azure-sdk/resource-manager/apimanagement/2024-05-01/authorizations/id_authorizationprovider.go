package authorizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AuthorizationProviderId{})
}

var _ resourceids.ResourceId = &AuthorizationProviderId{}

// AuthorizationProviderId is a struct representing the Resource ID for a Authorization Provider
type AuthorizationProviderId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ServiceName             string
	AuthorizationProviderId string
}

// NewAuthorizationProviderID returns a new AuthorizationProviderId struct
func NewAuthorizationProviderID(subscriptionId string, resourceGroupName string, serviceName string, authorizationProviderId string) AuthorizationProviderId {
	return AuthorizationProviderId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ServiceName:             serviceName,
		AuthorizationProviderId: authorizationProviderId,
	}
}

// ParseAuthorizationProviderID parses 'input' into a AuthorizationProviderId
func ParseAuthorizationProviderID(input string) (*AuthorizationProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthorizationProviderIDInsensitively parses 'input' case-insensitively into a AuthorizationProviderId
// note: this method should only be used for API response data and not user input
func ParseAuthorizationProviderIDInsensitively(input string) (*AuthorizationProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthorizationProviderId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AuthorizationProviderId, ok = input.Parsed["authorizationProviderId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationProviderId", input)
	}

	return nil
}

// ValidateAuthorizationProviderID checks that 'input' can be parsed as a Authorization Provider ID
func ValidateAuthorizationProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthorizationProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Authorization Provider ID
func (id AuthorizationProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/authorizationProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.AuthorizationProviderId)
}

// Segments returns a slice of Resource ID Segments which comprise this Authorization Provider ID
func (id AuthorizationProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticAuthorizationProviders", "authorizationProviders", "authorizationProviders"),
		resourceids.UserSpecifiedSegment("authorizationProviderId", "authorizationProviderId"),
	}
}

// String returns a human-readable description of this Authorization Provider ID
func (id AuthorizationProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Authorization Provider: %q", id.AuthorizationProviderId),
	}
	return fmt.Sprintf("Authorization Provider (%s)", strings.Join(components, "\n"))
}
