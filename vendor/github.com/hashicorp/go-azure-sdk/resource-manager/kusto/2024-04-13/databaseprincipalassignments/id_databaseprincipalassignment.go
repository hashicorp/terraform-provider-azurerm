package databaseprincipalassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabasePrincipalAssignmentId{})
}

var _ resourceids.ResourceId = &DatabasePrincipalAssignmentId{}

// DatabasePrincipalAssignmentId is a struct representing the Resource ID for a Database Principal Assignment
type DatabasePrincipalAssignmentId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ClusterName             string
	DatabaseName            string
	PrincipalAssignmentName string
}

// NewDatabasePrincipalAssignmentID returns a new DatabasePrincipalAssignmentId struct
func NewDatabasePrincipalAssignmentID(subscriptionId string, resourceGroupName string, clusterName string, databaseName string, principalAssignmentName string) DatabasePrincipalAssignmentId {
	return DatabasePrincipalAssignmentId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ClusterName:             clusterName,
		DatabaseName:            databaseName,
		PrincipalAssignmentName: principalAssignmentName,
	}
}

// ParseDatabasePrincipalAssignmentID parses 'input' into a DatabasePrincipalAssignmentId
func ParseDatabasePrincipalAssignmentID(input string) (*DatabasePrincipalAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabasePrincipalAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabasePrincipalAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabasePrincipalAssignmentIDInsensitively parses 'input' case-insensitively into a DatabasePrincipalAssignmentId
// note: this method should only be used for API response data and not user input
func ParseDatabasePrincipalAssignmentIDInsensitively(input string) (*DatabasePrincipalAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabasePrincipalAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabasePrincipalAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabasePrincipalAssignmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PrincipalAssignmentName, ok = input.Parsed["principalAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "principalAssignmentName", input)
	}

	return nil
}

// ValidateDatabasePrincipalAssignmentID checks that 'input' can be parsed as a Database Principal Assignment ID
func ValidateDatabasePrincipalAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabasePrincipalAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Principal Assignment ID
func (id DatabasePrincipalAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/clusters/%s/databases/%s/principalAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.DatabaseName, id.PrincipalAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Principal Assignment ID
func (id DatabasePrincipalAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKusto", "Microsoft.Kusto", "Microsoft.Kusto"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseName"),
		resourceids.StaticSegment("staticPrincipalAssignments", "principalAssignments", "principalAssignments"),
		resourceids.UserSpecifiedSegment("principalAssignmentName", "principalAssignmentName"),
	}
}

// String returns a human-readable description of this Database Principal Assignment ID
func (id DatabasePrincipalAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Principal Assignment Name: %q", id.PrincipalAssignmentName),
	}
	return fmt.Sprintf("Database Principal Assignment (%s)", strings.Join(components, "\n"))
}
