package openidconnectprovider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = OpenidConnectProviderId{}

// OpenidConnectProviderId is a struct representing the Resource ID for a Openid Connect Provider
type OpenidConnectProviderId struct {
	SubscriptionId            string
	ResourceGroupName         string
	ServiceName               string
	OpenidConnectProviderName string
}

// NewOpenidConnectProviderID returns a new OpenidConnectProviderId struct
func NewOpenidConnectProviderID(subscriptionId string, resourceGroupName string, serviceName string, openidConnectProviderName string) OpenidConnectProviderId {
	return OpenidConnectProviderId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		ServiceName:               serviceName,
		OpenidConnectProviderName: openidConnectProviderName,
	}
}

// ParseOpenidConnectProviderID parses 'input' into a OpenidConnectProviderId
func ParseOpenidConnectProviderID(input string) (*OpenidConnectProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(OpenidConnectProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OpenidConnectProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.OpenidConnectProviderName, ok = parsed.Parsed["openidConnectProviderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "openidConnectProviderName", *parsed)
	}

	return &id, nil
}

// ParseOpenidConnectProviderIDInsensitively parses 'input' case-insensitively into a OpenidConnectProviderId
// note: this method should only be used for API response data and not user input
func ParseOpenidConnectProviderIDInsensitively(input string) (*OpenidConnectProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(OpenidConnectProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OpenidConnectProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.OpenidConnectProviderName, ok = parsed.Parsed["openidConnectProviderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "openidConnectProviderName", *parsed)
	}

	return &id, nil
}

// ValidateOpenidConnectProviderID checks that 'input' can be parsed as a Openid Connect Provider ID
func ValidateOpenidConnectProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOpenidConnectProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Openid Connect Provider ID
func (id OpenidConnectProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/openidConnectProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.OpenidConnectProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Openid Connect Provider ID
func (id OpenidConnectProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticOpenidConnectProviders", "openidConnectProviders", "openidConnectProviders"),
		resourceids.UserSpecifiedSegment("openidConnectProviderName", "openidConnectProviderValue"),
	}
}

// String returns a human-readable description of this Openid Connect Provider ID
func (id OpenidConnectProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Openid Connect Provider Name: %q", id.OpenidConnectProviderName),
	}
	return fmt.Sprintf("Openid Connect Provider (%s)", strings.Join(components, "\n"))
}
