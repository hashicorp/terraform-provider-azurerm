package fluxconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FluxConfigurationId{}

// FluxConfigurationId is a struct representing the Resource ID for a Flux Configuration
type FluxConfigurationId struct {
	SubscriptionId        string
	ResourceGroupName     string
	ProviderName          string
	ClusterResourceName   string
	ClusterName           string
	FluxConfigurationName string
}

// NewFluxConfigurationID returns a new FluxConfigurationId struct
func NewFluxConfigurationID(subscriptionId string, resourceGroupName string, providerName string, clusterResourceName string, clusterName string, fluxConfigurationName string) FluxConfigurationId {
	return FluxConfigurationId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		ProviderName:          providerName,
		ClusterResourceName:   clusterResourceName,
		ClusterName:           clusterName,
		FluxConfigurationName: fluxConfigurationName,
	}
}

// ParseFluxConfigurationID parses 'input' into a FluxConfigurationId
func ParseFluxConfigurationID(input string) (*FluxConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(FluxConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FluxConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	if id.ClusterResourceName, ok = parsed.Parsed["clusterResourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterResourceName", *parsed)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterName", *parsed)
	}

	if id.FluxConfigurationName, ok = parsed.Parsed["fluxConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluxConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseFluxConfigurationIDInsensitively parses 'input' case-insensitively into a FluxConfigurationId
// note: this method should only be used for API response data and not user input
func ParseFluxConfigurationIDInsensitively(input string) (*FluxConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(FluxConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FluxConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	if id.ClusterResourceName, ok = parsed.Parsed["clusterResourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterResourceName", *parsed)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterName", *parsed)
	}

	if id.FluxConfigurationName, ok = parsed.Parsed["fluxConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluxConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateFluxConfigurationID checks that 'input' can be parsed as a Flux Configuration ID
func ValidateFluxConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFluxConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Flux Configuration ID
func (id FluxConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProviderName, id.ClusterResourceName, id.ClusterName, id.FluxConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Flux Configuration ID
func (id FluxConfigurationId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticFluxConfigurations", "fluxConfigurations", "fluxConfigurations"),
		resourceids.UserSpecifiedSegment("fluxConfigurationName", "fluxConfigurationValue"),
	}
}

// String returns a human-readable description of this Flux Configuration ID
func (id FluxConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Cluster Resource Name: %q", id.ClusterResourceName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Flux Configuration Name: %q", id.FluxConfigurationName),
	}
	return fmt.Sprintf("Flux Configuration (%s)", strings.Join(components, "\n"))
}
