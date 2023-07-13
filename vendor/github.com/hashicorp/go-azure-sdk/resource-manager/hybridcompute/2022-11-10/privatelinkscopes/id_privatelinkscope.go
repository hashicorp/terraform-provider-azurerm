package privatelinkscopes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrivateLinkScopeId{}

// PrivateLinkScopeId is a struct representing the Resource ID for a Private Link Scope
type PrivateLinkScopeId struct {
	SubscriptionId     string
	LocationName       string
	PrivateLinkScopeId string
}

// NewPrivateLinkScopeID returns a new PrivateLinkScopeId struct
func NewPrivateLinkScopeID(subscriptionId string, locationName string, privateLinkScopeId string) PrivateLinkScopeId {
	return PrivateLinkScopeId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		PrivateLinkScopeId: privateLinkScopeId,
	}
}

// ParsePrivateLinkScopeID parses 'input' into a PrivateLinkScopeId
func ParsePrivateLinkScopeID(input string) (*PrivateLinkScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateLinkScopeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateLinkScopeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.PrivateLinkScopeId, ok = parsed.Parsed["privateLinkScopeId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateLinkScopeId", *parsed)
	}

	return &id, nil
}

// ParsePrivateLinkScopeIDInsensitively parses 'input' case-insensitively into a PrivateLinkScopeId
// note: this method should only be used for API response data and not user input
func ParsePrivateLinkScopeIDInsensitively(input string) (*PrivateLinkScopeId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateLinkScopeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateLinkScopeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.PrivateLinkScopeId, ok = parsed.Parsed["privateLinkScopeId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateLinkScopeId", *parsed)
	}

	return &id, nil
}

// ValidatePrivateLinkScopeID checks that 'input' can be parsed as a Private Link Scope ID
func ValidatePrivateLinkScopeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateLinkScopeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Link Scope ID
func (id PrivateLinkScopeId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.HybridCompute/locations/%s/privateLinkScopes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.PrivateLinkScopeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Link Scope ID
func (id PrivateLinkScopeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticPrivateLinkScopes", "privateLinkScopes", "privateLinkScopes"),
		resourceids.UserSpecifiedSegment("privateLinkScopeId", "privateLinkScopeIdValue"),
	}
}

// String returns a human-readable description of this Private Link Scope ID
func (id PrivateLinkScopeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Private Link Scope: %q", id.PrivateLinkScopeId),
	}
	return fmt.Sprintf("Private Link Scope (%s)", strings.Join(components, "\n"))
}
