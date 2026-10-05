package privatelinkscopedresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ScopedResourceId{})
}

var _ resourceids.ResourceId = &ScopedResourceId{}

// ScopedResourceId is a struct representing the Resource ID for a Scoped Resource
type ScopedResourceId struct {
	SubscriptionId       string
	ResourceGroupName    string
	PrivateLinkScopeName string
	ScopedResourceName   string
}

// NewScopedResourceID returns a new ScopedResourceId struct
func NewScopedResourceID(subscriptionId string, resourceGroupName string, privateLinkScopeName string, scopedResourceName string) ScopedResourceId {
	return ScopedResourceId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		PrivateLinkScopeName: privateLinkScopeName,
		ScopedResourceName:   scopedResourceName,
	}
}

// ParseScopedResourceID parses 'input' into a ScopedResourceId
func ParseScopedResourceID(input string) (*ScopedResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedResourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseScopedResourceIDInsensitively parses 'input' case-insensitively into a ScopedResourceId
// note: this method should only be used for API response data and not user input
func ParseScopedResourceIDInsensitively(input string) (*ScopedResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ScopedResourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ScopedResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ScopedResourceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateLinkScopeName, ok = input.Parsed["privateLinkScopeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateLinkScopeName", input)
	}

	if id.ScopedResourceName, ok = input.Parsed["scopedResourceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "scopedResourceName", input)
	}

	return nil
}

// ValidateScopedResourceID checks that 'input' can be parsed as a Scoped Resource ID
func ValidateScopedResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopedResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scoped Resource ID
func (id ScopedResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/privateLinkScopes/%s/scopedResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateLinkScopeName, id.ScopedResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scoped Resource ID
func (id ScopedResourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticPrivateLinkScopes", "privateLinkScopes", "privateLinkScopes"),
		resourceids.UserSpecifiedSegment("privateLinkScopeName", "privateLinkScopeName"),
		resourceids.StaticSegment("staticScopedResources", "scopedResources", "scopedResources"),
		resourceids.UserSpecifiedSegment("scopedResourceName", "scopedResourceName"),
	}
}

// String returns a human-readable description of this Scoped Resource ID
func (id ScopedResourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Link Scope Name: %q", id.PrivateLinkScopeName),
		fmt.Sprintf("Scoped Resource Name: %q", id.ScopedResourceName),
	}
	return fmt.Sprintf("Scoped Resource (%s)", strings.Join(components, "\n"))
}
