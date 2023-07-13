package managedcassandras

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CassandraClusterId{}

// CassandraClusterId is a struct representing the Resource ID for a Cassandra Cluster
type CassandraClusterId struct {
	SubscriptionId       string
	ResourceGroupName    string
	CassandraClusterName string
}

// NewCassandraClusterID returns a new CassandraClusterId struct
func NewCassandraClusterID(subscriptionId string, resourceGroupName string, cassandraClusterName string) CassandraClusterId {
	return CassandraClusterId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		CassandraClusterName: cassandraClusterName,
	}
}

// ParseCassandraClusterID parses 'input' into a CassandraClusterId
func ParseCassandraClusterID(input string) (*CassandraClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(CassandraClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CassandraClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CassandraClusterName, ok = parsed.Parsed["cassandraClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cassandraClusterName", *parsed)
	}

	return &id, nil
}

// ParseCassandraClusterIDInsensitively parses 'input' case-insensitively into a CassandraClusterId
// note: this method should only be used for API response data and not user input
func ParseCassandraClusterIDInsensitively(input string) (*CassandraClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(CassandraClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CassandraClusterId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CassandraClusterName, ok = parsed.Parsed["cassandraClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cassandraClusterName", *parsed)
	}

	return &id, nil
}

// ValidateCassandraClusterID checks that 'input' can be parsed as a Cassandra Cluster ID
func ValidateCassandraClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCassandraClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cassandra Cluster ID
func (id CassandraClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/cassandraClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CassandraClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cassandra Cluster ID
func (id CassandraClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticCassandraClusters", "cassandraClusters", "cassandraClusters"),
		resourceids.UserSpecifiedSegment("cassandraClusterName", "cassandraClusterValue"),
	}
}

// String returns a human-readable description of this Cassandra Cluster ID
func (id CassandraClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cassandra Cluster Name: %q", id.CassandraClusterName),
	}
	return fmt.Sprintf("Cassandra Cluster (%s)", strings.Join(components, "\n"))
}
