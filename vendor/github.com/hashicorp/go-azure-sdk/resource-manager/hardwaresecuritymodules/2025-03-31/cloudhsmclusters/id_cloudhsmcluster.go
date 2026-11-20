package cloudhsmclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CloudHsmClusterId{})
}

var _ resourceids.ResourceId = &CloudHsmClusterId{}

// CloudHsmClusterId is a struct representing the Resource ID for a Cloud Hsm Cluster
type CloudHsmClusterId struct {
	SubscriptionId      string
	ResourceGroupName   string
	CloudHsmClusterName string
}

// NewCloudHsmClusterID returns a new CloudHsmClusterId struct
func NewCloudHsmClusterID(subscriptionId string, resourceGroupName string, cloudHsmClusterName string) CloudHsmClusterId {
	return CloudHsmClusterId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		CloudHsmClusterName: cloudHsmClusterName,
	}
}

// ParseCloudHsmClusterID parses 'input' into a CloudHsmClusterId
func ParseCloudHsmClusterID(input string) (*CloudHsmClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudHsmClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudHsmClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudHsmClusterIDInsensitively parses 'input' case-insensitively into a CloudHsmClusterId
// note: this method should only be used for API response data and not user input
func ParseCloudHsmClusterIDInsensitively(input string) (*CloudHsmClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudHsmClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudHsmClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudHsmClusterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudHsmClusterName, ok = input.Parsed["cloudHsmClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudHsmClusterName", input)
	}

	return nil
}

// ValidateCloudHsmClusterID checks that 'input' can be parsed as a Cloud Hsm Cluster ID
func ValidateCloudHsmClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudHsmClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud Hsm Cluster ID
func (id CloudHsmClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudHsmClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Hsm Cluster ID
func (id CloudHsmClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHardwareSecurityModules", "Microsoft.HardwareSecurityModules", "Microsoft.HardwareSecurityModules"),
		resourceids.StaticSegment("staticCloudHsmClusters", "cloudHsmClusters", "cloudHsmClusters"),
		resourceids.UserSpecifiedSegment("cloudHsmClusterName", "cloudHsmClusterName"),
	}
}

// String returns a human-readable description of this Cloud Hsm Cluster ID
func (id CloudHsmClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Hsm Cluster Name: %q", id.CloudHsmClusterName),
	}
	return fmt.Sprintf("Cloud Hsm Cluster (%s)", strings.Join(components, "\n"))
}
