// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &HDInsightClusterId{}

// HDInsightClusterId is a struct representing the Resource ID for a HDInsight Cluster
type HDInsightClusterId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
}

// NewHDInsightClusterID returns a new HDInsightClusterId struct
func NewHDInsightClusterID(subscriptionId string, resourceGroupName string, clusterName string) HDInsightClusterId {
	return HDInsightClusterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
	}
}

// ParseHDInsightClusterID parses 'input' into a HDInsightClusterId
func ParseHDInsightClusterID(input string) (*HDInsightClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HDInsightClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HDInsightClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHDInsightClusterIDInsensitively parses 'input' case-insensitively into a HDInsightClusterId
// note: this method should only be used for API response data and not user input
func ParseHDInsightClusterIDInsensitively(input string) (*HDInsightClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HDInsightClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HDInsightClusterId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HDInsightClusterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	return nil
}

// ValidateHDInsightClusterID checks that 'input' can be parsed as a HDInsight Cluster ID
func ValidateHDInsightClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHDInsightClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted HDInsight Cluster ID
func (id HDInsightClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HDInsight/clusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this HDInsight Cluster ID
func (id HDInsightClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHDInsight", "Microsoft.HDInsight", "Microsoft.HDInsight"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
	}
}

// String returns a human-readable description of this HDInsight Cluster ID
func (id HDInsightClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
	}
	return fmt.Sprintf("HDInsight Cluster (%s)", strings.Join(components, "\n"))
}
