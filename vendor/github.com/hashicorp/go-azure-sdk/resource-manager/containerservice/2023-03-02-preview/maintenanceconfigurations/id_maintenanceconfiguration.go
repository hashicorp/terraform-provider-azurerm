package maintenanceconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = MaintenanceConfigurationId{}

// MaintenanceConfigurationId is a struct representing the Resource ID for a Maintenance Configuration
type MaintenanceConfigurationId struct {
	SubscriptionId               string
	ResourceGroupName            string
	ManagedClusterName           string
	MaintenanceConfigurationName string
}

// NewMaintenanceConfigurationID returns a new MaintenanceConfigurationId struct
func NewMaintenanceConfigurationID(subscriptionId string, resourceGroupName string, managedClusterName string, maintenanceConfigurationName string) MaintenanceConfigurationId {
	return MaintenanceConfigurationId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		ManagedClusterName:           managedClusterName,
		MaintenanceConfigurationName: maintenanceConfigurationName,
	}
}

// ParseMaintenanceConfigurationID parses 'input' into a MaintenanceConfigurationId
func ParseMaintenanceConfigurationID(input string) (*MaintenanceConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(MaintenanceConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MaintenanceConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	if id.MaintenanceConfigurationName, ok = parsed.Parsed["maintenanceConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "maintenanceConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseMaintenanceConfigurationIDInsensitively parses 'input' case-insensitively into a MaintenanceConfigurationId
// note: this method should only be used for API response data and not user input
func ParseMaintenanceConfigurationIDInsensitively(input string) (*MaintenanceConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(MaintenanceConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MaintenanceConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedClusterName, ok = parsed.Parsed["managedClusterName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", *parsed)
	}

	if id.MaintenanceConfigurationName, ok = parsed.Parsed["maintenanceConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "maintenanceConfigurationName", *parsed)
	}

	return &id, nil
}

// ValidateMaintenanceConfigurationID checks that 'input' can be parsed as a Maintenance Configuration ID
func ValidateMaintenanceConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMaintenanceConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Maintenance Configuration ID
func (id MaintenanceConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/maintenanceConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, id.MaintenanceConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Maintenance Configuration ID
func (id MaintenanceConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterValue"),
		resourceids.StaticSegment("staticMaintenanceConfigurations", "maintenanceConfigurations", "maintenanceConfigurations"),
		resourceids.UserSpecifiedSegment("maintenanceConfigurationName", "maintenanceConfigurationValue"),
	}
}

// String returns a human-readable description of this Maintenance Configuration ID
func (id MaintenanceConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
		fmt.Sprintf("Maintenance Configuration Name: %q", id.MaintenanceConfigurationName),
	}
	return fmt.Sprintf("Maintenance Configuration (%s)", strings.Join(components, "\n"))
}
