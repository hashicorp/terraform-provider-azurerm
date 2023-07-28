package componentversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RegistryComponentVersionId{}

// RegistryComponentVersionId is a struct representing the Resource ID for a Registry Component Version
type RegistryComponentVersionId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	ComponentName     string
	VersionName       string
}

// NewRegistryComponentVersionID returns a new RegistryComponentVersionId struct
func NewRegistryComponentVersionID(subscriptionId string, resourceGroupName string, registryName string, componentName string, versionName string) RegistryComponentVersionId {
	return RegistryComponentVersionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		ComponentName:     componentName,
		VersionName:       versionName,
	}
}

// ParseRegistryComponentVersionID parses 'input' into a RegistryComponentVersionId
func ParseRegistryComponentVersionID(input string) (*RegistryComponentVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryComponentVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryComponentVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ComponentName, ok = parsed.Parsed["componentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "componentName", *parsed)
	}

	if id.VersionName, ok = parsed.Parsed["versionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionName", *parsed)
	}

	return &id, nil
}

// ParseRegistryComponentVersionIDInsensitively parses 'input' case-insensitively into a RegistryComponentVersionId
// note: this method should only be used for API response data and not user input
func ParseRegistryComponentVersionIDInsensitively(input string) (*RegistryComponentVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryComponentVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryComponentVersionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.ComponentName, ok = parsed.Parsed["componentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "componentName", *parsed)
	}

	if id.VersionName, ok = parsed.Parsed["versionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "versionName", *parsed)
	}

	return &id, nil
}

// ValidateRegistryComponentVersionID checks that 'input' can be parsed as a Registry Component Version ID
func ValidateRegistryComponentVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRegistryComponentVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Registry Component Version ID
func (id RegistryComponentVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/registries/%s/components/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.ComponentName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Registry Component Version ID
func (id RegistryComponentVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticComponents", "components", "components"),
		resourceids.UserSpecifiedSegment("componentName", "componentValue"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionValue"),
	}
}

// String returns a human-readable description of this Registry Component Version ID
func (id RegistryComponentVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Component Name: %q", id.ComponentName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Registry Component Version (%s)", strings.Join(components, "\n"))
}
