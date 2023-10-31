package identityprovider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = IdentityProviderId{}

// IdentityProviderId is a struct representing the Resource ID for a Identity Provider
type IdentityProviderId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ServiceName          string
	IdentityProviderName IdentityProviderType
}

// NewIdentityProviderID returns a new IdentityProviderId struct
func NewIdentityProviderID(subscriptionId string, resourceGroupName string, serviceName string, identityProviderName IdentityProviderType) IdentityProviderId {
	return IdentityProviderId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ServiceName:          serviceName,
		IdentityProviderName: identityProviderName,
	}
}

// ParseIdentityProviderID parses 'input' into a IdentityProviderId
func ParseIdentityProviderID(input string) (*IdentityProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(IdentityProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IdentityProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if v, ok := parsed.Parsed["identityProviderName"]; true {
		if !ok {
			return nil, resourceids.NewSegmentNotSpecifiedError(id, "identityProviderName", *parsed)
		}

		identityProviderName, err := parseIdentityProviderType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.IdentityProviderName = *identityProviderName
	}

	return &id, nil
}

// ParseIdentityProviderIDInsensitively parses 'input' case-insensitively into a IdentityProviderId
// note: this method should only be used for API response data and not user input
func ParseIdentityProviderIDInsensitively(input string) (*IdentityProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(IdentityProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IdentityProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if v, ok := parsed.Parsed["identityProviderName"]; true {
		if !ok {
			return nil, resourceids.NewSegmentNotSpecifiedError(id, "identityProviderName", *parsed)
		}

		identityProviderName, err := parseIdentityProviderType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.IdentityProviderName = *identityProviderName
	}

	return &id, nil
}

// ValidateIdentityProviderID checks that 'input' can be parsed as a Identity Provider ID
func ValidateIdentityProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIdentityProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Identity Provider ID
func (id IdentityProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/identityProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, string(id.IdentityProviderName))
}

// Segments returns a slice of Resource ID Segments which comprise this Identity Provider ID
func (id IdentityProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticIdentityProviders", "identityProviders", "identityProviders"),
		resourceids.ConstantSegment("identityProviderName", PossibleValuesForIdentityProviderType(), "aad"),
	}
}

// String returns a human-readable description of this Identity Provider ID
func (id IdentityProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Identity Provider Name: %q", string(id.IdentityProviderName)),
	}
	return fmt.Sprintf("Identity Provider (%s)", strings.Join(components, "\n"))
}
