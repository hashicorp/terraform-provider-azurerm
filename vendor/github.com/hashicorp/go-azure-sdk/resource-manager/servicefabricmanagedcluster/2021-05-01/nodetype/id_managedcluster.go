package nodetype

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ManagedClusterId{})
}

var _ resourceids.ResourceId = &ManagedClusterId{}

// ManagedClusterId is a struct representing the Resource ID for a Managed Cluster
type ManagedClusterId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ManagedClusterName string
}

// NewManagedClusterID returns a new ManagedClusterId struct
func NewManagedClusterID(subscriptionId string, resourceGroupName string, managedClusterName string) ManagedClusterId {
	return ManagedClusterId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ManagedClusterName: managedClusterName,
	}
}

// ParseManagedClusterID parses 'input' into a ManagedClusterId
func ParseManagedClusterID(input string) (*ManagedClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseManagedClusterIDInsensitively parses 'input' case-insensitively into a ManagedClusterId
// note: this method should only be used for API response data and not user input
func ParseManagedClusterIDInsensitively(input string) (*ManagedClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ManagedClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ManagedClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ManagedClusterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedClusterName, ok = input.Parsed["managedClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", input)
	}

	return nil
}

// ValidateManagedClusterID checks that 'input' can be parsed as a Managed Cluster ID
func ValidateManagedClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Cluster ID
func (id ManagedClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceFabric/managedClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Cluster ID
func (id ManagedClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceFabric", "Microsoft.ServiceFabric", "Microsoft.ServiceFabric"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterName"),
	}
}

// String returns a human-readable description of this Managed Cluster ID
func (id ManagedClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
	}
	return fmt.Sprintf("Managed Cluster (%s)", strings.Join(components, "\n"))
}
