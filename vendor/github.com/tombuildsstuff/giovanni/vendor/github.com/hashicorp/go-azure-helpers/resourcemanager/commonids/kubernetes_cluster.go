// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = KubernetesClusterId{}

// KubernetesClusterId is a struct representing the Resource ID for a Kubernetes Cluster
type KubernetesClusterId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ManagedClusterName string
}

// NewKubernetesClusterID returns a new KubernetesClusterId struct
func NewKubernetesClusterID(subscriptionId string, resourceGroupName string, managedClusterName string) KubernetesClusterId {
	return KubernetesClusterId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ManagedClusterName: managedClusterName,
	}
}

// ParseKubernetesClusterID parses 'input' into a KubernetesClusterId
func ParseKubernetesClusterID(input string) (*KubernetesClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(KubernetesClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KubernetesClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	return &id, nil
}

// ParseKubernetesClusterIdInsensitively parses 'input' case-insensitively into a KubernetesClusterId
// note: this method should only be used for API response data and not user input
func ParseKubernetesClusterIDInsensitively(input string) (*KubernetesClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(KubernetesClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KubernetesClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	return &id, nil
}

// ValidateKubernetesClusterID checks that 'input' can be parsed as a Kubernetes Cluster ID
func ValidateKubernetesClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKubernetesClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Kubernetes Cluster ID
func (id KubernetesClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Kubernetes Cluster ID
func (id KubernetesClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterValue"),
	}
}

// String returns a human-readable description of this Kubernetes Cluster ID
func (id KubernetesClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Kubernetes Cluster Name: %q", id.ManagedClusterName),
	}
	return fmt.Sprintf("Kubernetes Cluster (%s)", strings.Join(components, "\n"))
}
