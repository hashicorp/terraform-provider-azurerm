package maintenanceconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &MaintenanceConfigurationId{}

// MaintenanceConfigurationId is a struct representing the Resource ID for a Maintenance Configuration
type MaintenanceConfigurationId struct {
	SubscriptionId               string
	ResourceGroupName            string
	MaintenanceConfigurationName string
}

// NewMaintenanceConfigurationID returns a new MaintenanceConfigurationId struct
func NewMaintenanceConfigurationID(subscriptionId string, resourceGroupName string, maintenanceConfigurationName string) MaintenanceConfigurationId {
	return MaintenanceConfigurationId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		MaintenanceConfigurationName: maintenanceConfigurationName,
	}
}

// ParseMaintenanceConfigurationID parses 'input' into a MaintenanceConfigurationId
func ParseMaintenanceConfigurationID(input string) (*MaintenanceConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MaintenanceConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MaintenanceConfigurationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMaintenanceConfigurationIDInsensitively parses 'input' case-insensitively into a MaintenanceConfigurationId
// note: this method should only be used for API response data and not user input
func ParseMaintenanceConfigurationIDInsensitively(input string) (*MaintenanceConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MaintenanceConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MaintenanceConfigurationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MaintenanceConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MaintenanceConfigurationName, ok = input.Parsed["maintenanceConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "maintenanceConfigurationName", input)
	}

	return nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Maintenance/maintenanceConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MaintenanceConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Maintenance Configuration ID
func (id MaintenanceConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaintenance", "Microsoft.Maintenance", "Microsoft.Maintenance"),
		resourceids.StaticSegment("staticMaintenanceConfigurations", "maintenanceConfigurations", "maintenanceConfigurations"),
		resourceids.UserSpecifiedSegment("maintenanceConfigurationName", "maintenanceConfigurationValue"),
	}
}

// String returns a human-readable description of this Maintenance Configuration ID
func (id MaintenanceConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Maintenance Configuration Name: %q", id.MaintenanceConfigurationName),
	}
	return fmt.Sprintf("Maintenance Configuration (%s)", strings.Join(components, "\n"))
}
