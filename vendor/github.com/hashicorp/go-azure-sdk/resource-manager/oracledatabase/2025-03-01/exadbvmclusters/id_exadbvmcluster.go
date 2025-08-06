package exadbvmclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExadbVMClusterId{})
}

var _ resourceids.ResourceId = &ExadbVMClusterId{}

// ExadbVMClusterId is a struct representing the Resource ID for a Exadb V M Cluster
type ExadbVMClusterId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ExadbVmClusterName string
}

// NewExadbVMClusterID returns a new ExadbVMClusterId struct
func NewExadbVMClusterID(subscriptionId string, resourceGroupName string, exadbVmClusterName string) ExadbVMClusterId {
	return ExadbVMClusterId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ExadbVmClusterName: exadbVmClusterName,
	}
}

// ParseExadbVMClusterID parses 'input' into a ExadbVMClusterId
func ParseExadbVMClusterID(input string) (*ExadbVMClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExadbVMClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExadbVMClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExadbVMClusterIDInsensitively parses 'input' case-insensitively into a ExadbVMClusterId
// note: this method should only be used for API response data and not user input
func ParseExadbVMClusterIDInsensitively(input string) (*ExadbVMClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExadbVMClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExadbVMClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExadbVMClusterId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateExadbVMClusterID checks that 'input' can be parsed as a Exadb V M Cluster ID
func ValidateExadbVMClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExadbVMClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Exadb V M Cluster ID
func (id ExadbVMClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/exadbVmClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExadbVmClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Exadb V M Cluster ID
func (id ExadbVMClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticExadbVmClusters", "exadbVmClusters", "exadbVmClusters"),
		resourceids.UserSpecifiedSegment("exadbVmClusterName", "exadbVmClusterName"),
	}
}

// String returns a human-readable description of this Exadb V M Cluster ID
func (id ExadbVMClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Exadb Vm Cluster Name: %q", id.ExadbVmClusterName),
	}
	return fmt.Sprintf("Exadb V M Cluster (%s)", strings.Join(components, "\n"))
}
