// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = KustoClusterId{}

// KustoClusterId is a struct representing the Resource ID for a Kusto Cluster
type KustoClusterId struct {
	SubscriptionId    string
	ResourceGroupName string
	KustoClusterName  string
}

// NewKustoClusterID returns a new KustoClusterId struct
func NewKustoClusterID(subscriptionId string, resourceGroupName string, kustoClusterName string) KustoClusterId {
	return KustoClusterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		KustoClusterName:  kustoClusterName,
	}
}

// ParseKustoClusterID parses 'input' into a KustoClusterId
func ParseKustoClusterID(input string) (*KustoClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(KustoClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KustoClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.KustoClusterName, ok = parsed.Parsed["kustoClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "kustoClusterName", *parsed)
	}

	return &id, nil
}

// ParseKustoClusterIDInsensitively parses 'input' case-insensitively into a KustoClusterId
// note: this method should only be used for API response data and not user input
func ParseKustoClusterIDInsensitively(input string) (*KustoClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(KustoClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KustoClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.KustoClusterName, ok = parsed.Parsed["kustoClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "kustoClusterName", *parsed)
	}

	return &id, nil
}

// ValidateKustoClusterID checks that 'input' can be parsed as a Kusto Cluster ID
func ValidateKustoClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKustoClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Kusto Cluster ID
func (id KustoClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.KustoClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Kusto Cluster ID
func (id KustoClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("kustoClusterName", "clusterValue"),
	}
}

// String returns a human-readable description of this Kusto Cluster ID
func (id KustoClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Kusto Cluster Name: %q", id.KustoClusterName),
	}
	return fmt.Sprintf("Cluster (%s)", strings.Join(components, "\n"))
}
