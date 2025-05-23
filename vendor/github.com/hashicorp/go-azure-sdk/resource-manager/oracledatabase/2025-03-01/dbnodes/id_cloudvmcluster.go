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
	recaser.RegisterResourceId(&CloudVMClusterId{})
}

var _ resourceids.ResourceId = &CloudVMClusterId{}

// CloudVMClusterId is a struct representing the Resource ID for a Cloud V M Cluster
type CloudVMClusterId struct {
	SubscriptionId     string
	ResourceGroupName  string
	CloudVmClusterName string
}

// NewCloudVMClusterID returns a new CloudVMClusterId struct
func NewCloudVMClusterID(subscriptionId string, resourceGroupName string, cloudVmClusterName string) CloudVMClusterId {
	return CloudVMClusterId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		CloudVmClusterName: cloudVmClusterName,
	}
}

// ParseCloudVMClusterID parses 'input' into a CloudVMClusterId
func ParseCloudVMClusterID(input string) (*CloudVMClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudVMClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudVMClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudVMClusterIDInsensitively parses 'input' case-insensitively into a CloudVMClusterId
// note: this method should only be used for API response data and not user input
func ParseCloudVMClusterIDInsensitively(input string) (*CloudVMClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudVMClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudVMClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudVMClusterId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateCloudVMClusterID checks that 'input' can be parsed as a Cloud V M Cluster ID
func ValidateCloudVMClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudVMClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud V M Cluster ID
func (id CloudVMClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/cloudVmClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudVmClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud V M Cluster ID
func (id CloudVMClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticCloudVmClusters", "cloudVmClusters", "cloudVmClusters"),
		resourceids.UserSpecifiedSegment("cloudVmClusterName", "cloudVmClusterName"),
	}
}

// String returns a human-readable description of this Cloud V M Cluster ID
func (id CloudVMClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Vm Cluster Name: %q", id.CloudVmClusterName),
	}
	return fmt.Sprintf("Cloud V M Cluster (%s)", strings.Join(components, "\n"))
}
