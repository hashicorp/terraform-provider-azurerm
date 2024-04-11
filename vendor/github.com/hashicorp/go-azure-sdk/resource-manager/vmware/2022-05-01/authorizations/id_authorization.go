package authorizations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &AuthorizationId{}

// AuthorizationId is a struct representing the Resource ID for a Authorization
type AuthorizationId struct {
	SubscriptionId    string
	ResourceGroupName string
	PrivateCloudName  string
	AuthorizationName string
}

// NewAuthorizationID returns a new AuthorizationId struct
func NewAuthorizationID(subscriptionId string, resourceGroupName string, privateCloudName string, authorizationName string) AuthorizationId {
	return AuthorizationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PrivateCloudName:  privateCloudName,
		AuthorizationName: authorizationName,
	}
}

// ParseAuthorizationID parses 'input' into a AuthorizationId
func ParseAuthorizationID(input string) (*AuthorizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthorizationIDInsensitively parses 'input' case-insensitively into a AuthorizationId
// note: this method should only be used for API response data and not user input
func ParseAuthorizationIDInsensitively(input string) (*AuthorizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthorizationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthorizationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthorizationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateCloudName, ok = input.Parsed["privateCloudName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateCloudName", input)
	}

	if id.AuthorizationName, ok = input.Parsed["authorizationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authorizationName", input)
	}

	return nil
}

// ValidateAuthorizationID checks that 'input' can be parsed as a Authorization ID
func ValidateAuthorizationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthorizationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Authorization ID
func (id AuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName, id.AuthorizationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Authorization ID
func (id AuthorizationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAVS", "Microsoft.AVS", "Microsoft.AVS"),
		resourceids.StaticSegment("staticPrivateClouds", "privateClouds", "privateClouds"),
		resourceids.UserSpecifiedSegment("privateCloudName", "privateCloudValue"),
		resourceids.StaticSegment("staticAuthorizations", "authorizations", "authorizations"),
		resourceids.UserSpecifiedSegment("authorizationName", "authorizationValue"),
	}
}

// String returns a human-readable description of this Authorization ID
func (id AuthorizationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Cloud Name: %q", id.PrivateCloudName),
		fmt.Sprintf("Authorization Name: %q", id.AuthorizationName),
	}
	return fmt.Sprintf("Authorization (%s)", strings.Join(components, "\n"))
}
