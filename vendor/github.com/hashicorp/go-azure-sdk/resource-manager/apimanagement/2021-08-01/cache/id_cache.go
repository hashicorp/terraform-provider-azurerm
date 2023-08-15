package cache

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CacheId{}

// CacheId is a struct representing the Resource ID for a Cache
type CacheId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	CacheId           string
}

// NewCacheID returns a new CacheId struct
func NewCacheID(subscriptionId string, resourceGroupName string, serviceName string, cacheId string) CacheId {
	return CacheId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		CacheId:           cacheId,
	}
}

// ParseCacheID parses 'input' into a CacheId
func ParseCacheID(input string) (*CacheId, error) {
	parser := resourceids.NewParserFromResourceIdType(CacheId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CacheId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.CacheId, ok = parsed.Parsed["cacheId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cacheId", *parsed)
	}

	return &id, nil
}

// ParseCacheIDInsensitively parses 'input' case-insensitively into a CacheId
// note: this method should only be used for API response data and not user input
func ParseCacheIDInsensitively(input string) (*CacheId, error) {
	parser := resourceids.NewParserFromResourceIdType(CacheId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CacheId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.CacheId, ok = parsed.Parsed["cacheId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cacheId", *parsed)
	}

	return &id, nil
}

// ValidateCacheID checks that 'input' can be parsed as a Cache ID
func ValidateCacheID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCacheID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cache ID
func (id CacheId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/caches/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.CacheId)
}

// Segments returns a slice of Resource ID Segments which comprise this Cache ID
func (id CacheId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticCaches", "caches", "caches"),
		resourceids.UserSpecifiedSegment("cacheId", "cacheIdValue"),
	}
}

// String returns a human-readable description of this Cache ID
func (id CacheId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Cache: %q", id.CacheId),
	}
	return fmt.Sprintf("Cache (%s)", strings.Join(components, "\n"))
}
