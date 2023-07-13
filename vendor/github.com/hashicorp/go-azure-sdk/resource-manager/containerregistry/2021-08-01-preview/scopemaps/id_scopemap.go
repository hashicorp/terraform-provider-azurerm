package scopemaps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ScopeMapId{}

// ScopeMapId is a struct representing the Resource ID for a Scope Map
type ScopeMapId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	ScopeMapName      string
}

// NewScopeMapID returns a new ScopeMapId struct
func NewScopeMapID(subscriptionId string, resourceGroupName string, registryName string, scopeMapName string) ScopeMapId {
	return ScopeMapId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		ScopeMapName:      scopeMapName,
	}
}

// ParseScopeMapID parses 'input' into a ScopeMapId
func ParseScopeMapID(input string) (*ScopeMapId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopeMapId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopeMapId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ScopeMapName, ok = parsed.Parsed["scopeMapName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scopeMapName", *parsed)
	}

	return &id, nil
}

// ParseScopeMapIDInsensitively parses 'input' case-insensitively into a ScopeMapId
// note: this method should only be used for API response data and not user input
func ParseScopeMapIDInsensitively(input string) (*ScopeMapId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScopeMapId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScopeMapId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ScopeMapName, ok = parsed.Parsed["scopeMapName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "scopeMapName", *parsed)
	}

	return &id, nil
}

// ValidateScopeMapID checks that 'input' can be parsed as a Scope Map ID
func ValidateScopeMapID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScopeMapID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scope Map ID
func (id ScopeMapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/scopeMaps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ScopeMapName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scope Map ID
func (id ScopeMapId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticScopeMaps", "scopeMaps", "scopeMaps"),
		resourceids.UserSpecifiedSegment("scopeMapName", "scopeMapValue"),
	}
}

// String returns a human-readable description of this Scope Map ID
func (id ScopeMapId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Scope Map Name: %q", id.ScopeMapName),
	}
	return fmt.Sprintf("Scope Map (%s)", strings.Join(components, "\n"))
}
