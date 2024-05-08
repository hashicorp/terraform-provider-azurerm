// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &KustoDatabaseId{}

// KustoDatabaseId is a struct representing the Resource ID for a Kusto Database
type KustoDatabaseId struct {
	SubscriptionId    string
	ResourceGroupName string
	KustoClusterName  string
	KustoDatabaseName string
}

// NewKustoDatabaseID returns a new KustoDatabaseId struct
func NewKustoDatabaseID(subscriptionId string, resourceGroupName string, kustoClusterName string, kustoDatabaseName string) KustoDatabaseId {
	return KustoDatabaseId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		KustoClusterName:  kustoClusterName,
		KustoDatabaseName: kustoDatabaseName,
	}
}

// ParseKustoDatabaseID parses 'input' into a KustoDatabaseId
func ParseKustoDatabaseID(input string) (*KustoDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KustoDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KustoDatabaseId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKustoDatabaseIDInsensitively parses 'input' case-insensitively into a KustoDatabaseId
// note: this method should only be used for API response data and not user input
func ParseKustoDatabaseIDInsensitively(input string) (*KustoDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KustoDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KustoDatabaseId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KustoDatabaseId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.KustoClusterName, ok = input.Parsed["kustoClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "kustoClusterName", input)
	}

	if id.KustoDatabaseName, ok = input.Parsed["kustoDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "kustoDatabaseName", input)
	}

	return nil
}

// ValidateKustoDatabaseID checks that 'input' can be parsed as a KustoDatabase ID
func ValidateKustoDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKustoDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Kusto Database ID
func (id KustoDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/databases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.KustoClusterName, id.KustoDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Kusto Database ID
func (id KustoDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("kustoClusterName", "clusterValue"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("kustoDatabaseName", "databaseValue"),
	}
}

// String returns a human-readable description of this Kusto Database ID
func (id KustoDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Kusto Cluster Name: %q", id.KustoClusterName),
		fmt.Sprintf("Kusto Database Name: %q", id.KustoDatabaseName),
	}
	return fmt.Sprintf("Kusto Database (%s)", strings.Join(components, "\n"))
}
