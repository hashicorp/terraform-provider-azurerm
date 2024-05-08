package hdinsights

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ClusterPoolId{})
}

var _ resourceids.ResourceId = &ClusterPoolId{}

// ClusterPoolId is a struct representing the Resource ID for a Cluster Pool
type ClusterPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterPoolName   string
}

// NewClusterPoolID returns a new ClusterPoolId struct
func NewClusterPoolID(subscriptionId string, resourceGroupName string, clusterPoolName string) ClusterPoolId {
	return ClusterPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterPoolName:   clusterPoolName,
	}
}

// ParseClusterPoolID parses 'input' into a ClusterPoolId
func ParseClusterPoolID(input string) (*ClusterPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClusterPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClusterPoolId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseClusterPoolIDInsensitively parses 'input' case-insensitively into a ClusterPoolId
// note: this method should only be used for API response data and not user input
func ParseClusterPoolIDInsensitively(input string) (*ClusterPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClusterPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClusterPoolId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ClusterPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterPoolName, ok = input.Parsed["clusterPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterPoolName", input)
	}

	return nil
}

// ValidateClusterPoolID checks that 'input' can be parsed as a Cluster Pool ID
func ValidateClusterPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseClusterPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cluster Pool ID
func (id ClusterPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HDInsight/clusterPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cluster Pool ID
func (id ClusterPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHDInsight", "Microsoft.HDInsight", "Microsoft.HDInsight"),
		resourceids.StaticSegment("staticClusterPools", "clusterPools", "clusterPools"),
		resourceids.UserSpecifiedSegment("clusterPoolName", "clusterPoolValue"),
	}
}

// String returns a human-readable description of this Cluster Pool ID
func (id ClusterPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Pool Name: %q", id.ClusterPoolName),
	}
	return fmt.Sprintf("Cluster Pool (%s)", strings.Join(components, "\n"))
}
