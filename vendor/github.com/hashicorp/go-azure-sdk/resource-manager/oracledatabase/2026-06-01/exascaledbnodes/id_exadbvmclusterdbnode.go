package exascaledbnodes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExadbVMClusterDbNodeId{})
}

var _ resourceids.ResourceId = &ExadbVMClusterDbNodeId{}

// ExadbVMClusterDbNodeId is a struct representing the Resource ID for a Exadb VM Cluster Db Node
type ExadbVMClusterDbNodeId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ExadbVmClusterName string
	DbNodeName         string
}

// NewExadbVMClusterDbNodeID returns a new ExadbVMClusterDbNodeId struct
func NewExadbVMClusterDbNodeID(subscriptionId string, resourceGroupName string, exadbVmClusterName string, dbNodeName string) ExadbVMClusterDbNodeId {
	return ExadbVMClusterDbNodeId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ExadbVmClusterName: exadbVmClusterName,
		DbNodeName:         dbNodeName,
	}
}

// ParseExadbVMClusterDbNodeID parses 'input' into a ExadbVMClusterDbNodeId
func ParseExadbVMClusterDbNodeID(input string) (*ExadbVMClusterDbNodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExadbVMClusterDbNodeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExadbVMClusterDbNodeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExadbVMClusterDbNodeIDInsensitively parses 'input' case-insensitively into a ExadbVMClusterDbNodeId
// note: this method should only be used for API response data and not user input
func ParseExadbVMClusterDbNodeIDInsensitively(input string) (*ExadbVMClusterDbNodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExadbVMClusterDbNodeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExadbVMClusterDbNodeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExadbVMClusterDbNodeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExadbVmClusterName, ok = input.Parsed["exadbVmClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "exadbVmClusterName", input)
	}

	if id.DbNodeName, ok = input.Parsed["dbNodeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dbNodeName", input)
	}

	return nil
}

// ValidateExadbVMClusterDbNodeID checks that 'input' can be parsed as a Exadb VM Cluster Db Node ID
func ValidateExadbVMClusterDbNodeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExadbVMClusterDbNodeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Exadb VM Cluster Db Node ID
func (id ExadbVMClusterDbNodeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/exadbVmClusters/%s/dbNodes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExadbVmClusterName, id.DbNodeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Exadb VM Cluster Db Node ID
func (id ExadbVMClusterDbNodeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticExadbVmClusters", "exadbVmClusters", "exadbVmClusters"),
		resourceids.UserSpecifiedSegment("exadbVmClusterName", "exadbVmClusterName"),
		resourceids.StaticSegment("staticDbNodes", "dbNodes", "dbNodes"),
		resourceids.UserSpecifiedSegment("dbNodeName", "dbNodeName"),
	}
}

// String returns a human-readable description of this Exadb VM Cluster Db Node ID
func (id ExadbVMClusterDbNodeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Exadb Vm Cluster Name: %q", id.ExadbVmClusterName),
		fmt.Sprintf("Db Node Name: %q", id.DbNodeName),
	}
	return fmt.Sprintf("Exadb VM Cluster Db Node (%s)", strings.Join(components, "\n"))
}
