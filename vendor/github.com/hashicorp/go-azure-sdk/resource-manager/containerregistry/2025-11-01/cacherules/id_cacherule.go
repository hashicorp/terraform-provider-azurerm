package cacherules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CacheRuleId{})
}

var _ resourceids.ResourceId = &CacheRuleId{}

// CacheRuleId is a struct representing the Resource ID for a Cache Rule
type CacheRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	CacheRuleName     string
}

// NewCacheRuleID returns a new CacheRuleId struct
func NewCacheRuleID(subscriptionId string, resourceGroupName string, registryName string, cacheRuleName string) CacheRuleId {
	return CacheRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		CacheRuleName:     cacheRuleName,
	}
}

// ParseCacheRuleID parses 'input' into a CacheRuleId
func ParseCacheRuleID(input string) (*CacheRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CacheRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CacheRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCacheRuleIDInsensitively parses 'input' case-insensitively into a CacheRuleId
// note: this method should only be used for API response data and not user input
func ParseCacheRuleIDInsensitively(input string) (*CacheRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CacheRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CacheRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CacheRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.CacheRuleName, ok = input.Parsed["cacheRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cacheRuleName", input)
	}

	return nil
}

// ValidateCacheRuleID checks that 'input' can be parsed as a Cache Rule ID
func ValidateCacheRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCacheRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cache Rule ID
func (id CacheRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/cacheRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.CacheRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cache Rule ID
func (id CacheRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticCacheRules", "cacheRules", "cacheRules"),
		resourceids.UserSpecifiedSegment("cacheRuleName", "cacheRuleName"),
	}
}

// String returns a human-readable description of this Cache Rule ID
func (id CacheRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Cache Rule Name: %q", id.CacheRuleName),
	}
	return fmt.Sprintf("Cache Rule (%s)", strings.Join(components, "\n"))
}
