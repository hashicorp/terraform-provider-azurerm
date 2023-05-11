package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderPrivateLinkScopeId{}

// ProviderPrivateLinkScopeId is a struct representing the Resource ID for a Provider Private Link Scope
type ProviderPrivateLinkScopeId struct {
	SubscriptionId       string
	ResourceGroupName    string
	PrivateLinkScopeName string
}

// NewProviderPrivateLinkScopeID returns a new ProviderPrivateLinkScopeId struct
func NewProviderPrivateLinkScopeID(subscriptionId string, resourceGroupName string, privateLinkScopeName string) ProviderPrivateLinkScopeId {
	return ProviderPrivateLinkScopeId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		PrivateLinkScopeName: privateLinkScopeName,
	}
}

// ParseProviderPrivateLinkScopeID parses 'input' into a ProviderPrivateLinkScopeId
func ParseProviderPrivateLinkScopeID(input string) (*ProviderPrivateLinkScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderPrivateLinkScopeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderPrivateLinkScopeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateLinkScopeName, ok = parsed.Parsed["privateLinkScopeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateLinkScopeName", *parsed)
	}

	return &id, nil
}

// ParseProviderPrivateLinkScopeIDInsensitively parses 'input' case-insensitively into a ProviderPrivateLinkScopeId
// note: this method should only be used for API response data and not user input
func ParseProviderPrivateLinkScopeIDInsensitively(input string) (*ProviderPrivateLinkScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderPrivateLinkScopeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderPrivateLinkScopeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateLinkScopeName, ok = parsed.Parsed["privateLinkScopeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateLinkScopeName", *parsed)
	}

	return &id, nil
}

// ValidateProviderPrivateLinkScopeID checks that 'input' can be parsed as a Provider Private Link Scope ID
func ValidateProviderPrivateLinkScopeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderPrivateLinkScopeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Private Link Scope ID
func (id ProviderPrivateLinkScopeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/privateLinkScopes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateLinkScopeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Private Link Scope ID
func (id ProviderPrivateLinkScopeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticPrivateLinkScopes", "privateLinkScopes", "privateLinkScopes"),
		resourceids.UserSpecifiedSegment("privateLinkScopeName", "privateLinkScopeValue"),
	}
}

// String returns a human-readable description of this Provider Private Link Scope ID
func (id ProviderPrivateLinkScopeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Link Scope Name: %q", id.PrivateLinkScopeName),
	}
	return fmt.Sprintf("Provider Private Link Scope (%s)", strings.Join(components, "\n"))
}
