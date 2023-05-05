package extensions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ExtensionId{}

// ExtensionId is a struct representing the Resource ID for a Extension
type ExtensionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ProviderName        string
	ClusterResourceName string
	ClusterName         string
	ExtensionName       string
}

// NewExtensionID returns a new ExtensionId struct
func NewExtensionID(subscriptionId string, resourceGroupName string, providerName string, clusterResourceName string, clusterName string, extensionName string) ExtensionId {
	return ExtensionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ProviderName:        providerName,
		ClusterResourceName: clusterResourceName,
		ClusterName:         clusterName,
		ExtensionName:       extensionName,
	}
}

// ParseExtensionID parses 'input' into a ExtensionId
func ParseExtensionID(input string) (*ExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExtensionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExtensionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ClusterResourceName, ok = parsed.Parsed["clusterResourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterResourceName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.ExtensionName, ok = parsed.Parsed["extensionName"]; !ok {
		return nil, fmt.Errorf("the segment 'extensionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseExtensionIDInsensitively parses 'input' case-insensitively into a ExtensionId
// note: this method should only be used for API response data and not user input
func ParseExtensionIDInsensitively(input string) (*ExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ExtensionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ExtensionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ClusterResourceName, ok = parsed.Parsed["clusterResourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterResourceName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.ExtensionName, ok = parsed.Parsed["extensionName"]; !ok {
		return nil, fmt.Errorf("the segment 'extensionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateExtensionID checks that 'input' can be parsed as a Extension ID
func ValidateExtensionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExtensionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Extension ID
func (id ExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.KubernetesConfiguration/extensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProviderName, id.ClusterResourceName, id.ClusterName, id.ExtensionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Extension ID
func (id ExtensionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
		resourceids.UserSpecifiedSegment("clusterResourceName", "clusterResourceValue"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKubernetesConfiguration", "Microsoft.KubernetesConfiguration", "Microsoft.KubernetesConfiguration"),
		resourceids.StaticSegment("staticExtensions", "extensions", "extensions"),
		resourceids.UserSpecifiedSegment("extensionName", "extensionValue"),
	}
}

// String returns a human-readable description of this Extension ID
func (id ExtensionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Cluster Resource Name: %q", id.ClusterResourceName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Extension Name: %q", id.ExtensionName),
	}
	return fmt.Sprintf("Extension (%s)", strings.Join(components, "\n"))
}
