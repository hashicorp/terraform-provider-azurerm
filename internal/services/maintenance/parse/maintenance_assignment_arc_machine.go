// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &MaintenanceAssignmentArcMachineId{}

// MaintenanceAssignmentArcMachineId restricts the SDK's generic scope to an Arc Machine.
// The generic scoped ID also accepts VM and Dedicated Host assignments. Explicit Arc
// segments validate the parent during import and expose a hierarchical Resource Identity.
type MaintenanceAssignmentArcMachineId struct {
	SubscriptionId              string
	ResourceGroupName           string
	MachineName                 string
	ConfigurationAssignmentName string
}

func NewMaintenanceAssignmentArcMachineID(subscriptionId, resourceGroupName, machineName, configurationAssignmentName string) MaintenanceAssignmentArcMachineId {
	return MaintenanceAssignmentArcMachineId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		MachineName:                 machineName,
		ConfigurationAssignmentName: configurationAssignmentName,
	}
}

func MaintenanceAssignmentArcMachineID(input string) (*MaintenanceAssignmentArcMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MaintenanceAssignmentArcMachineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MaintenanceAssignmentArcMachineId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// MaintenanceAssignmentArcMachineIDInsensitively is for API response data, not user input.
func MaintenanceAssignmentArcMachineIDInsensitively(input string) (*MaintenanceAssignmentArcMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MaintenanceAssignmentArcMachineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MaintenanceAssignmentArcMachineId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MaintenanceAssignmentArcMachineId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MachineName, ok = input.Parsed["machineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "machineName", input)
	}

	if id.ConfigurationAssignmentName, ok = input.Parsed["configurationAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationAssignmentName", input)
	}

	return nil
}

func (id MaintenanceAssignmentArcMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/providers/Microsoft.Maintenance/configurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MachineName, id.ConfigurationAssignmentName)
}

func (MaintenanceAssignmentArcMachineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineName"),
		resourceids.StaticSegment("staticMaintenanceProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaintenance", "Microsoft.Maintenance", "Microsoft.Maintenance"),
		resourceids.StaticSegment("staticConfigurationAssignments", "configurationAssignments", "configurationAssignments"),
		resourceids.UserSpecifiedSegment("configurationAssignmentName", "configurationAssignmentName"),
	}
}

func (id MaintenanceAssignmentArcMachineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
		fmt.Sprintf("Configuration Assignment Name: %q", id.ConfigurationAssignmentName),
	}
	return fmt.Sprintf("Maintenance Assignment Arc Machine (%s)", strings.Join(components, "\n"))
}
