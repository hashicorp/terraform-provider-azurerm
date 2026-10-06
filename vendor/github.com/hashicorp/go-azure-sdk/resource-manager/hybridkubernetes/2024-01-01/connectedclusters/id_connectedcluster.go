package connectedclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectedClusterId{})
}

var _ resourceids.ResourceId = &ConnectedClusterId{}

// ConnectedClusterId is a struct representing the Resource ID for a Connected Cluster
type ConnectedClusterId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ConnectedClusterName string
}

// NewConnectedClusterID returns a new ConnectedClusterId struct
func NewConnectedClusterID(subscriptionId string, resourceGroupName string, connectedClusterName string) ConnectedClusterId {
	return ConnectedClusterId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ConnectedClusterName: connectedClusterName,
	}
}

// ParseConnectedClusterID parses 'input' into a ConnectedClusterId
func ParseConnectedClusterID(input string) (*ConnectedClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectedClusterIDInsensitively parses 'input' case-insensitively into a ConnectedClusterId
// note: this method should only be used for API response data and not user input
func ParseConnectedClusterIDInsensitively(input string) (*ConnectedClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectedClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectedClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectedClusterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ConnectedClusterName, ok = input.Parsed["connectedClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectedClusterName", input)
	}

	return nil
}

// ValidateConnectedClusterID checks that 'input' can be parsed as a Connected Cluster ID
func ValidateConnectedClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectedClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connected Cluster ID
func (id ConnectedClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kubernetes/connectedClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConnectedClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connected Cluster ID
func (id ConnectedClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKubernetes", "Microsoft.Kubernetes", "Microsoft.Kubernetes"),
		resourceids.StaticSegment("staticConnectedClusters", "connectedClusters", "connectedClusters"),
		resourceids.UserSpecifiedSegment("connectedClusterName", "connectedClusterName"),
	}
}

// String returns a human-readable description of this Connected Cluster ID
func (id ConnectedClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Connected Cluster Name: %q", id.ConnectedClusterName),
	}
	return fmt.Sprintf("Connected Cluster (%s)", strings.Join(components, "\n"))
}
