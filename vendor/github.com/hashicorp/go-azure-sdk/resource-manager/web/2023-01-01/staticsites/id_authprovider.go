package staticsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AuthProviderId{})
}

var _ resourceids.ResourceId = &AuthProviderId{}

// AuthProviderId is a struct representing the Resource ID for a Auth Provider
type AuthProviderId struct {
	SubscriptionId    string
	ResourceGroupName string
	StaticSiteName    string
	AuthProviderName  string
}

// NewAuthProviderID returns a new AuthProviderId struct
func NewAuthProviderID(subscriptionId string, resourceGroupName string, staticSiteName string, authProviderName string) AuthProviderId {
	return AuthProviderId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StaticSiteName:    staticSiteName,
		AuthProviderName:  authProviderName,
	}
}

// ParseAuthProviderID parses 'input' into a AuthProviderId
func ParseAuthProviderID(input string) (*AuthProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthProviderIDInsensitively parses 'input' case-insensitively into a AuthProviderId
// note: this method should only be used for API response data and not user input
func ParseAuthProviderIDInsensitively(input string) (*AuthProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthProviderId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthProviderId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StaticSiteName, ok = input.Parsed["staticSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "staticSiteName", input)
	}

	if id.AuthProviderName, ok = input.Parsed["authProviderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authProviderName", input)
	}

	return nil
}

// ValidateAuthProviderID checks that 'input' can be parsed as a Auth Provider ID
func ValidateAuthProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Auth Provider ID
func (id AuthProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/staticSites/%s/authProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName, id.AuthProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Auth Provider ID
func (id AuthProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticStaticSites", "staticSites", "staticSites"),
		resourceids.UserSpecifiedSegment("staticSiteName", "staticSiteName"),
		resourceids.StaticSegment("staticAuthProviders", "authProviders", "authProviders"),
		resourceids.UserSpecifiedSegment("authProviderName", "authProviderName"),
	}
}

// String returns a human-readable description of this Auth Provider ID
func (id AuthProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Static Site Name: %q", id.StaticSiteName),
		fmt.Sprintf("Auth Provider Name: %q", id.AuthProviderName),
	}
	return fmt.Sprintf("Auth Provider (%s)", strings.Join(components, "\n"))
}
