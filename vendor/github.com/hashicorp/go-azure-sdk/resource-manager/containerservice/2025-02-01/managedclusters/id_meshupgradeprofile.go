package managedclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MeshUpgradeProfileId{})
}

var _ resourceids.ResourceId = &MeshUpgradeProfileId{}

// MeshUpgradeProfileId is a struct representing the Resource ID for a Mesh Upgrade Profile
type MeshUpgradeProfileId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedClusterName     string
	MeshUpgradeProfileName string
}

// NewMeshUpgradeProfileID returns a new MeshUpgradeProfileId struct
func NewMeshUpgradeProfileID(subscriptionId string, resourceGroupName string, managedClusterName string, meshUpgradeProfileName string) MeshUpgradeProfileId {
	return MeshUpgradeProfileId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedClusterName:     managedClusterName,
		MeshUpgradeProfileName: meshUpgradeProfileName,
	}
}

// ParseMeshUpgradeProfileID parses 'input' into a MeshUpgradeProfileId
func ParseMeshUpgradeProfileID(input string) (*MeshUpgradeProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MeshUpgradeProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MeshUpgradeProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMeshUpgradeProfileIDInsensitively parses 'input' case-insensitively into a MeshUpgradeProfileId
// note: this method should only be used for API response data and not user input
func ParseMeshUpgradeProfileIDInsensitively(input string) (*MeshUpgradeProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MeshUpgradeProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MeshUpgradeProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MeshUpgradeProfileId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.MeshUpgradeProfileName, ok = input.Parsed["meshUpgradeProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "meshUpgradeProfileName", input)
	}

	return nil
}

// ValidateMeshUpgradeProfileID checks that 'input' can be parsed as a Mesh Upgrade Profile ID
func ValidateMeshUpgradeProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMeshUpgradeProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mesh Upgrade Profile ID
func (id MeshUpgradeProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/meshUpgradeProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, id.MeshUpgradeProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mesh Upgrade Profile ID
func (id MeshUpgradeProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterName"),
		resourceids.StaticSegment("staticMeshUpgradeProfiles", "meshUpgradeProfiles", "meshUpgradeProfiles"),
		resourceids.UserSpecifiedSegment("meshUpgradeProfileName", "meshUpgradeProfileName"),
	}
}

// String returns a human-readable description of this Mesh Upgrade Profile ID
func (id MeshUpgradeProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
		fmt.Sprintf("Mesh Upgrade Profile Name: %q", id.MeshUpgradeProfileName),
	}
	return fmt.Sprintf("Mesh Upgrade Profile (%s)", strings.Join(components, "\n"))
}
