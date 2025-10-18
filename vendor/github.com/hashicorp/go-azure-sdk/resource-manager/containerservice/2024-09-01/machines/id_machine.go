package machines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MachineId{})
}

var _ resourceids.ResourceId = &MachineId{}

// MachineId is a struct representing the Resource ID for a Machine
type MachineId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ManagedClusterName string
	AgentPoolName      string
	MachineName        string
}

// NewMachineID returns a new MachineId struct
func NewMachineID(subscriptionId string, resourceGroupName string, managedClusterName string, agentPoolName string, machineName string) MachineId {
	return MachineId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ManagedClusterName: managedClusterName,
		AgentPoolName:      agentPoolName,
		MachineName:        machineName,
	}
}

// ParseMachineID parses 'input' into a MachineId
func ParseMachineID(input string) (*MachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MachineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MachineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMachineIDInsensitively parses 'input' case-insensitively into a MachineId
// note: this method should only be used for API response data and not user input
func ParseMachineIDInsensitively(input string) (*MachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MachineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MachineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MachineId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AgentPoolName, ok = input.Parsed["agentPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "agentPoolName", input)
	}

	if id.MachineName, ok = input.Parsed["machineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "machineName", input)
	}

	return nil
}

// ValidateMachineID checks that 'input' can be parsed as a Machine ID
func ValidateMachineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMachineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Machine ID
func (id MachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/agentPools/%s/machines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, id.AgentPoolName, id.MachineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Machine ID
func (id MachineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterName"),
		resourceids.StaticSegment("staticAgentPools", "agentPools", "agentPools"),
		resourceids.UserSpecifiedSegment("agentPoolName", "agentPoolName"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineName"),
	}
}

// String returns a human-readable description of this Machine ID
func (id MachineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
		fmt.Sprintf("Agent Pool Name: %q", id.AgentPoolName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
	}
	return fmt.Sprintf("Machine (%s)", strings.Join(components, "\n"))
}
