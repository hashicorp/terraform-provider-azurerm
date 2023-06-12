package managedclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AccessProfileId{}

// AccessProfileId is a struct representing the Resource ID for a Access Profile
type AccessProfileId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ManagedClusterName string
	AccessProfileName  string
}

// NewAccessProfileID returns a new AccessProfileId struct
func NewAccessProfileID(subscriptionId string, resourceGroupName string, managedClusterName string, accessProfileName string) AccessProfileId {
	return AccessProfileId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ManagedClusterName: managedClusterName,
		AccessProfileName:  accessProfileName,
	}
}

// ParseAccessProfileID parses 'input' into a AccessProfileId
func ParseAccessProfileID(input string) (*AccessProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	if id.AccessProfileName, ok = parsed.Parsed["accessProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accessProfileName", *parsed)
	}

	return &id, nil
}

// ParseAccessProfileIDInsensitively parses 'input' case-insensitively into a AccessProfileId
// note: this method should only be used for API response data and not user input
func ParseAccessProfileIDInsensitively(input string) (*AccessProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	if id.AccessProfileName, ok = parsed.Parsed["accessProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accessProfileName", *parsed)
	}

	return &id, nil
}

// ValidateAccessProfileID checks that 'input' can be parsed as a Access Profile ID
func ValidateAccessProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccessProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Access Profile ID
func (id AccessProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/accessProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, id.AccessProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Access Profile ID
func (id AccessProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterValue"),
		resourceids.StaticSegment("staticAccessProfiles", "accessProfiles", "accessProfiles"),
		resourceids.UserSpecifiedSegment("accessProfileName", "accessProfileValue"),
	}
}

// String returns a human-readable description of this Access Profile ID
func (id AccessProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
		fmt.Sprintf("Access Profile Name: %q", id.AccessProfileName),
	}
	return fmt.Sprintf("Access Profile (%s)", strings.Join(components, "\n"))
}
