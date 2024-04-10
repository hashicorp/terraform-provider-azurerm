package dataconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DataConnectionId{}

// DataConnectionId is a struct representing the Resource ID for a Data Connection
type DataConnectionId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ClusterName        string
	DatabaseName       string
	DataConnectionName string
}

// NewDataConnectionID returns a new DataConnectionId struct
func NewDataConnectionID(subscriptionId string, resourceGroupName string, clusterName string, databaseName string, dataConnectionName string) DataConnectionId {
	return DataConnectionId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ClusterName:        clusterName,
		DatabaseName:       databaseName,
		DataConnectionName: dataConnectionName,
	}
}

// ParseDataConnectionID parses 'input' into a DataConnectionId
func ParseDataConnectionID(input string) (*DataConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataConnectionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataConnectionIDInsensitively parses 'input' case-insensitively into a DataConnectionId
// note: this method should only be used for API response data and not user input
func ParseDataConnectionIDInsensitively(input string) (*DataConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataConnectionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataConnectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DatabaseName, ok = input.Parsed["databaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseName", input)
	}

	if id.DataConnectionName, ok = input.Parsed["dataConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataConnectionName", input)
	}

	return nil
}

// ValidateDataConnectionID checks that 'input' can be parsed as a Data Connection ID
func ValidateDataConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Connection ID
func (id DataConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/databases/%s/dataConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.DatabaseName, id.DataConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Connection ID
func (id DataConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseValue"),
		resourceids.StaticSegment("staticDataConnections", "dataConnections", "dataConnections"),
		resourceids.UserSpecifiedSegment("dataConnectionName", "dataConnectionValue"),
	}
}

// String returns a human-readable description of this Data Connection ID
func (id DataConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Data Connection Name: %q", id.DataConnectionName),
	}
	return fmt.Sprintf("Data Connection (%s)", strings.Join(components, "\n"))
}
