package expressroutelagauthorizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExpressRouteLagAuthorizationId{})
}

var _ resourceids.ResourceId = &ExpressRouteLagAuthorizationId{}

// ExpressRouteLagAuthorizationId is a struct representing the Resource ID for a Express Route Lag Authorization
type ExpressRouteLagAuthorizationId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ExpressRouteLagName string
	AuthorizationName   string
}

// NewExpressRouteLagAuthorizationID returns a new ExpressRouteLagAuthorizationId struct
func NewExpressRouteLagAuthorizationID(subscriptionId string, resourceGroupName string, expressRouteLagName string, authorizationName string) ExpressRouteLagAuthorizationId {
	return ExpressRouteLagAuthorizationId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ExpressRouteLagName: expressRouteLagName,
		AuthorizationName:   authorizationName,
	}
}

// ParseExpressRouteLagAuthorizationID parses 'input' into a ExpressRouteLagAuthorizationId
func ParseExpressRouteLagAuthorizationID(input string) (*ExpressRouteLagAuthorizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteLagAuthorizationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteLagAuthorizationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExpressRouteLagAuthorizationIDInsensitively parses 'input' case-insensitively into a ExpressRouteLagAuthorizationId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteLagAuthorizationIDInsensitively(input string) (*ExpressRouteLagAuthorizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteLagAuthorizationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteLagAuthorizationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExpressRouteLagAuthorizationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteLagName, ok = input.Parsed["expressRouteLagName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteLagName", input)
	}

	if id.AuthorizationName, ok = input.Parsed["authorizationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationName", input)
	}

	return nil
}

// ValidateExpressRouteLagAuthorizationID checks that 'input' can be parsed as a Express Route Lag Authorization ID
func ValidateExpressRouteLagAuthorizationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteLagAuthorizationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Lag Authorization ID
func (id ExpressRouteLagAuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteLags/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteLagName, id.AuthorizationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Lag Authorization ID
func (id ExpressRouteLagAuthorizationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteLags", "expressRouteLags", "expressRouteLags"),
		resourceids.UserSpecifiedSegment("expressRouteLagName", "expressRouteLagName"),
		resourceids.StaticSegment("staticAuthorizations", "authorizations", "authorizations"),
		resourceids.UserSpecifiedSegment("authorizationName", "authorizationName"),
	}
}

// String returns a human-readable description of this Express Route Lag Authorization ID
func (id ExpressRouteLagAuthorizationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Lag Name: %q", id.ExpressRouteLagName),
		fmt.Sprintf("Authorization Name: %q", id.AuthorizationName),
	}
	return fmt.Sprintf("Express Route Lag Authorization (%s)", strings.Join(components, "\n"))
}
