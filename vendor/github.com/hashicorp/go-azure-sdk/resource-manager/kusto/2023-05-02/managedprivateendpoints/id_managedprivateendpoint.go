package managedprivateendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ManagedPrivateEndpointId{}

// ManagedPrivateEndpointId is a struct representing the Resource ID for a Managed Private Endpoint
type ManagedPrivateEndpointId struct {
	SubscriptionId             string
	ResourceGroupName          string
	ClusterName                string
	ManagedPrivateEndpointName string
}

// NewManagedPrivateEndpointID returns a new ManagedPrivateEndpointId struct
func NewManagedPrivateEndpointID(subscriptionId string, resourceGroupName string, clusterName string, managedPrivateEndpointName string) ManagedPrivateEndpointId {
	return ManagedPrivateEndpointId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		ClusterName:                clusterName,
		ManagedPrivateEndpointName: managedPrivateEndpointName,
	}
}

// ParseManagedPrivateEndpointID parses 'input' into a ManagedPrivateEndpointId
func ParseManagedPrivateEndpointID(input string) (*ManagedPrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedPrivateEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedPrivateEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterName", *parsed)
	}

	if id.ManagedPrivateEndpointName, ok = parsed.Parsed["managedPrivateEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedPrivateEndpointName", *parsed)
	}

	return &id, nil
}

// ParseManagedPrivateEndpointIDInsensitively parses 'input' case-insensitively into a ManagedPrivateEndpointId
// note: this method should only be used for API response data and not user input
func ParseManagedPrivateEndpointIDInsensitively(input string) (*ManagedPrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedPrivateEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedPrivateEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clusterName", *parsed)
	}

	if id.ManagedPrivateEndpointName, ok = parsed.Parsed["managedPrivateEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedPrivateEndpointName", *parsed)
	}

	return &id, nil
}

// ValidateManagedPrivateEndpointID checks that 'input' can be parsed as a Managed Private Endpoint ID
func ValidateManagedPrivateEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedPrivateEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Private Endpoint ID
func (id ManagedPrivateEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/managedPrivateEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ManagedPrivateEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Private Endpoint ID
func (id ManagedPrivateEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticManagedPrivateEndpoints", "managedPrivateEndpoints", "managedPrivateEndpoints"),
		resourceids.UserSpecifiedSegment("managedPrivateEndpointName", "managedPrivateEndpointValue"),
	}
}

// String returns a human-readable description of this Managed Private Endpoint ID
func (id ManagedPrivateEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Managed Private Endpoint Name: %q", id.ManagedPrivateEndpointName),
	}
	return fmt.Sprintf("Managed Private Endpoint (%s)", strings.Join(components, "\n"))
}
