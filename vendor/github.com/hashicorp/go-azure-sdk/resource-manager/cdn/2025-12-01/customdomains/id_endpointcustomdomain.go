package customdomains

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EndpointCustomDomainId{})
}

var _ resourceids.ResourceId = &EndpointCustomDomainId{}

// EndpointCustomDomainId is a struct representing the Resource ID for a Endpoint Custom Domain
type EndpointCustomDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	EndpointName      string
	CustomDomainName  string
}

// NewEndpointCustomDomainID returns a new EndpointCustomDomainId struct
func NewEndpointCustomDomainID(subscriptionId string, resourceGroupName string, profileName string, endpointName string, customDomainName string) EndpointCustomDomainId {
	return EndpointCustomDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		EndpointName:      endpointName,
		CustomDomainName:  customDomainName,
	}
}

// ParseEndpointCustomDomainID parses 'input' into a EndpointCustomDomainId
func ParseEndpointCustomDomainID(input string) (*EndpointCustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EndpointCustomDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EndpointCustomDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEndpointCustomDomainIDInsensitively parses 'input' case-insensitively into a EndpointCustomDomainId
// note: this method should only be used for API response data and not user input
func ParseEndpointCustomDomainIDInsensitively(input string) (*EndpointCustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EndpointCustomDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EndpointCustomDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EndpointCustomDomainId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProfileName, ok = input.Parsed["profileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "profileName", input)
	}

	if id.EndpointName, ok = input.Parsed["endpointName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "endpointName", input)
	}

	if id.CustomDomainName, ok = input.Parsed["customDomainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "customDomainName", input)
	}

	return nil
}

// ValidateEndpointCustomDomainID checks that 'input' can be parsed as a Endpoint Custom Domain ID
func ValidateEndpointCustomDomainID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEndpointCustomDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Endpoint Custom Domain ID
func (id EndpointCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/endpoints/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName, id.CustomDomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Endpoint Custom Domain ID
func (id EndpointCustomDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticEndpoints", "endpoints", "endpoints"),
		resourceids.UserSpecifiedSegment("endpointName", "endpointName"),
		resourceids.StaticSegment("staticCustomDomains", "customDomains", "customDomains"),
		resourceids.UserSpecifiedSegment("customDomainName", "customDomainName"),
	}
}

// String returns a human-readable description of this Endpoint Custom Domain ID
func (id EndpointCustomDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Endpoint Name: %q", id.EndpointName),
		fmt.Sprintf("Custom Domain Name: %q", id.CustomDomainName),
	}
	return fmt.Sprintf("Endpoint Custom Domain (%s)", strings.Join(components, "\n"))
}
