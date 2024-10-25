package dbnodes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DbNodeId{})
}

var _ resourceids.ResourceId = &DbNodeId{}

// DbNodeId is a struct representing the Resource ID for a Db Node
type DbNodeId struct {
	SubscriptionId     string
	ResourceGroupName  string
	CloudVmClusterName string
	DbNodeName         string
}

// NewDbNodeID returns a new DbNodeId struct
func NewDbNodeID(subscriptionId string, resourceGroupName string, cloudVmClusterName string, dbNodeName string) DbNodeId {
	return DbNodeId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		CloudVmClusterName: cloudVmClusterName,
		DbNodeName:         dbNodeName,
	}
}

// ParseDbNodeID parses 'input' into a DbNodeId
func ParseDbNodeID(input string) (*DbNodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbNodeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbNodeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDbNodeIDInsensitively parses 'input' case-insensitively into a DbNodeId
// note: this method should only be used for API response data and not user input
func ParseDbNodeIDInsensitively(input string) (*DbNodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbNodeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbNodeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DbNodeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudVmClusterName, ok = input.Parsed["cloudVmClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudVmClusterName", input)
	}

	if id.DbNodeName, ok = input.Parsed["dbNodeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dbNodeName", input)
	}

	return nil
}

// ValidateDbNodeID checks that 'input' can be parsed as a Db Node ID
func ValidateDbNodeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDbNodeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Db Node ID
func (id DbNodeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/cloudVmClusters/%s/dbNodes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudVmClusterName, id.DbNodeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Db Node ID
func (id DbNodeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticCloudVmClusters", "cloudVmClusters", "cloudVmClusters"),
		resourceids.UserSpecifiedSegment("cloudVmClusterName", "cloudVmClusterName"),
		resourceids.StaticSegment("staticDbNodes", "dbNodes", "dbNodes"),
		resourceids.UserSpecifiedSegment("dbNodeName", "dbNodeName"),
	}
}

// String returns a human-readable description of this Db Node ID
func (id DbNodeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Vm Cluster Name: %q", id.CloudVmClusterName),
		fmt.Sprintf("Db Node Name: %q", id.DbNodeName),
	}
	return fmt.Sprintf("Db Node (%s)", strings.Join(components, "\n"))
}
