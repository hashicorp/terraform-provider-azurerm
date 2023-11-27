package expressrouteportauthorizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExpressRoutePortAuthorizationId{}

// ExpressRoutePortAuthorizationId is a struct representing the Resource ID for a Express Route Port Authorization
type ExpressRoutePortAuthorizationId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ExpressRoutePortName string
	AuthorizationName    string
}

// NewExpressRoutePortAuthorizationID returns a new ExpressRoutePortAuthorizationId struct
func NewExpressRoutePortAuthorizationID(subscriptionId string, resourceGroupName string, expressRoutePortName string, authorizationName string) ExpressRoutePortAuthorizationId {
	return ExpressRoutePortAuthorizationId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ExpressRoutePortName: expressRoutePortName,
		AuthorizationName:    authorizationName,
	}
}

// ParseExpressRoutePortAuthorizationID parses 'input' into a ExpressRoutePortAuthorizationId
func ParseExpressRoutePortAuthorizationID(input string) (*ExpressRoutePortAuthorizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRoutePortAuthorizationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRoutePortAuthorizationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRoutePortName, ok = parsed.Parsed["expressRoutePortName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortName", *parsed)
	}

	if id.AuthorizationName, ok = parsed.Parsed["authorizationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationName", *parsed)
	}

	return &id, nil
}

// ParseExpressRoutePortAuthorizationIDInsensitively parses 'input' case-insensitively into a ExpressRoutePortAuthorizationId
// note: this method should only be used for API response data and not user input
func ParseExpressRoutePortAuthorizationIDInsensitively(input string) (*ExpressRoutePortAuthorizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExpressRoutePortAuthorizationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExpressRoutePortAuthorizationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ExpressRoutePortName, ok = parsed.Parsed["expressRoutePortName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "expressRoutePortName", *parsed)
	}

	if id.AuthorizationName, ok = parsed.Parsed["authorizationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "authorizationName", *parsed)
	}

	return &id, nil
}

// ValidateExpressRoutePortAuthorizationID checks that 'input' can be parsed as a Express Route Port Authorization ID
func ValidateExpressRoutePortAuthorizationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRoutePortAuthorizationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Port Authorization ID
func (id ExpressRoutePortAuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRoutePorts/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRoutePortName, id.AuthorizationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Port Authorization ID
func (id ExpressRoutePortAuthorizationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRoutePorts", "expressRoutePorts", "expressRoutePorts"),
		resourceids.UserSpecifiedSegment("expressRoutePortName", "expressRoutePortValue"),
		resourceids.StaticSegment("staticAuthorizations", "authorizations", "authorizations"),
		resourceids.UserSpecifiedSegment("authorizationName", "authorizationValue"),
	}
}

// String returns a human-readable description of this Express Route Port Authorization ID
func (id ExpressRoutePortAuthorizationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Port Name: %q", id.ExpressRoutePortName),
		fmt.Sprintf("Authorization Name: %q", id.AuthorizationName),
	}
	return fmt.Sprintf("Express Route Port Authorization (%s)", strings.Join(components, "\n"))
}
