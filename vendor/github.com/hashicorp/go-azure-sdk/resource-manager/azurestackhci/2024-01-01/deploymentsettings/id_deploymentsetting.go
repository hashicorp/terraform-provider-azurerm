package deploymentsettings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DeploymentSettingId{})
}

var _ resourceids.ResourceId = &DeploymentSettingId{}

// DeploymentSettingId is a struct representing the Resource ID for a Deployment Setting
type DeploymentSettingId struct {
	SubscriptionId        string
	ResourceGroupName     string
	ClusterName           string
	DeploymentSettingName string
}

// NewDeploymentSettingID returns a new DeploymentSettingId struct
func NewDeploymentSettingID(subscriptionId string, resourceGroupName string, clusterName string, deploymentSettingName string) DeploymentSettingId {
	return DeploymentSettingId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		ClusterName:           clusterName,
		DeploymentSettingName: deploymentSettingName,
	}
}

// ParseDeploymentSettingID parses 'input' into a DeploymentSettingId
func ParseDeploymentSettingID(input string) (*DeploymentSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeploymentSettingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeploymentSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeploymentSettingIDInsensitively parses 'input' case-insensitively into a DeploymentSettingId
// note: this method should only be used for API response data and not user input
func ParseDeploymentSettingIDInsensitively(input string) (*DeploymentSettingId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeploymentSettingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeploymentSettingId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeploymentSettingId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DeploymentSettingName, ok = input.Parsed["deploymentSettingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deploymentSettingName", input)
	}

	return nil
}

// ValidateDeploymentSettingID checks that 'input' can be parsed as a Deployment Setting ID
func ValidateDeploymentSettingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeploymentSettingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deployment Setting ID
func (id DeploymentSettingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/deploymentSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.DeploymentSettingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deployment Setting ID
func (id DeploymentSettingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticDeploymentSettings", "deploymentSettings", "deploymentSettings"),
		resourceids.UserSpecifiedSegment("deploymentSettingName", "deploymentSettingName"),
	}
}

// String returns a human-readable description of this Deployment Setting ID
func (id DeploymentSettingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Deployment Setting Name: %q", id.DeploymentSettingName),
	}
	return fmt.Sprintf("Deployment Setting (%s)", strings.Join(components, "\n"))
}
