package maintenances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MaintenanceId{})
}

var _ resourceids.ResourceId = &MaintenanceId{}

// MaintenanceId is a struct representing the Resource ID for a Maintenance
type MaintenanceId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
	MaintenanceName    string
}

// NewMaintenanceID returns a new MaintenanceId struct
func NewMaintenanceID(subscriptionId string, resourceGroupName string, flexibleServerName string, maintenanceName string) MaintenanceId {
	return MaintenanceId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
		MaintenanceName:    maintenanceName,
	}
}

// ParseMaintenanceID parses 'input' into a MaintenanceId
func ParseMaintenanceID(input string) (*MaintenanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MaintenanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MaintenanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMaintenanceIDInsensitively parses 'input' case-insensitively into a MaintenanceId
// note: this method should only be used for API response data and not user input
func ParseMaintenanceIDInsensitively(input string) (*MaintenanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MaintenanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MaintenanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MaintenanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FlexibleServerName, ok = input.Parsed["flexibleServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexibleServerName", input)
	}

	if id.MaintenanceName, ok = input.Parsed["maintenanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "maintenanceName", input)
	}

	return nil
}

// ValidateMaintenanceID checks that 'input' can be parsed as a Maintenance ID
func ValidateMaintenanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMaintenanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Maintenance ID
func (id MaintenanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s/maintenances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName, id.MaintenanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Maintenance ID
func (id MaintenanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMySQL", "Microsoft.DBforMySQL", "Microsoft.DBforMySQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerName"),
		resourceids.StaticSegment("staticMaintenances", "maintenances", "maintenances"),
		resourceids.UserSpecifiedSegment("maintenanceName", "maintenanceName"),
	}
}

// String returns a human-readable description of this Maintenance ID
func (id MaintenanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
		fmt.Sprintf("Maintenance Name: %q", id.MaintenanceName),
	}
	return fmt.Sprintf("Maintenance (%s)", strings.Join(components, "\n"))
}
