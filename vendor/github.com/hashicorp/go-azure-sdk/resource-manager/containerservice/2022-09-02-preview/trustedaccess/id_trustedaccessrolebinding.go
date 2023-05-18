package trustedaccess

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TrustedAccessRoleBindingId{}

// TrustedAccessRoleBindingId is a struct representing the Resource ID for a Trusted Access Role Binding
type TrustedAccessRoleBindingId struct {
	SubscriptionId               string
	ResourceGroupName            string
	ManagedClusterName           string
	TrustedAccessRoleBindingName string
}

// NewTrustedAccessRoleBindingID returns a new TrustedAccessRoleBindingId struct
func NewTrustedAccessRoleBindingID(subscriptionId string, resourceGroupName string, managedClusterName string, trustedAccessRoleBindingName string) TrustedAccessRoleBindingId {
	return TrustedAccessRoleBindingId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		ManagedClusterName:           managedClusterName,
		TrustedAccessRoleBindingName: trustedAccessRoleBindingName,
	}
}

// ParseTrustedAccessRoleBindingID parses 'input' into a TrustedAccessRoleBindingId
func ParseTrustedAccessRoleBindingID(input string) (*TrustedAccessRoleBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrustedAccessRoleBindingId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrustedAccessRoleBindingId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	if id.TrustedAccessRoleBindingName, ok = parsed.Parsed["trustedAccessRoleBindingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trustedAccessRoleBindingName", *parsed)
	}

	return &id, nil
}

// ParseTrustedAccessRoleBindingIDInsensitively parses 'input' case-insensitively into a TrustedAccessRoleBindingId
// note: this method should only be used for API response data and not user input
func ParseTrustedAccessRoleBindingIDInsensitively(input string) (*TrustedAccessRoleBindingId, error) {
	parser := resourceids.NewParserFromResourceIdType(TrustedAccessRoleBindingId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TrustedAccessRoleBindingId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	if id.TrustedAccessRoleBindingName, ok = parsed.Parsed["trustedAccessRoleBindingName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trustedAccessRoleBindingName", *parsed)
	}

	return &id, nil
}

// ValidateTrustedAccessRoleBindingID checks that 'input' can be parsed as a Trusted Access Role Binding ID
func ValidateTrustedAccessRoleBindingID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTrustedAccessRoleBindingID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Trusted Access Role Binding ID
func (id TrustedAccessRoleBindingId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/trustedAccessRoleBindings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, id.TrustedAccessRoleBindingName)
}

// Segments returns a slice of Resource ID Segments which comprise this Trusted Access Role Binding ID
func (id TrustedAccessRoleBindingId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterValue"),
		resourceids.StaticSegment("staticTrustedAccessRoleBindings", "trustedAccessRoleBindings", "trustedAccessRoleBindings"),
		resourceids.UserSpecifiedSegment("trustedAccessRoleBindingName", "trustedAccessRoleBindingValue"),
	}
}

// String returns a human-readable description of this Trusted Access Role Binding ID
func (id TrustedAccessRoleBindingId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
		fmt.Sprintf("Trusted Access Role Binding Name: %q", id.TrustedAccessRoleBindingName),
	}
	return fmt.Sprintf("Trusted Access Role Binding (%s)", strings.Join(components, "\n"))
}
