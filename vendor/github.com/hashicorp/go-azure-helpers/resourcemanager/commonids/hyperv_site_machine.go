// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &HyperVSiteMachineId{}

// HyperVSiteMachineId is a struct representing the Resource ID for a Hyper V Site Machine
type HyperVSiteMachineId struct {
	SubscriptionId    string
	ResourceGroupName string
	HyperVSiteName    string
	MachineName       string
}

// NewHyperVSiteMachineID returns a new HyperVSiteMachineId struct
func NewHyperVSiteMachineID(subscriptionId string, resourceGroupName string, hyperVSiteName string, machineName string) HyperVSiteMachineId {
	return HyperVSiteMachineId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HyperVSiteName:    hyperVSiteName,
		MachineName:       machineName,
	}
}

// ParseHyperVSiteMachineID parses 'input' into a HyperVSiteMachineId
func ParseHyperVSiteMachineID(input string) (*HyperVSiteMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HyperVSiteMachineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HyperVSiteMachineId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHyperVSiteMachineIDInsensitively parses 'input' case-insensitively into a HyperVSiteMachineId
// note: this method should only be used for API response data and not user input
func ParseHyperVSiteMachineIDInsensitively(input string) (*HyperVSiteMachineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HyperVSiteMachineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HyperVSiteMachineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HyperVSiteMachineId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HyperVSiteName, ok = input.Parsed["hyperVSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hyperVSiteName", input)
	}

	if id.MachineName, ok = input.Parsed["machineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "machineName", input)
	}

	return nil
}

// ValidateHyperVSiteMachineID checks that 'input' can be parsed as a Hyper V Site Machine ID
func ValidateHyperVSiteMachineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHyperVSiteMachineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hyper V Site Machine ID
func (id HyperVSiteMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OffAzure/hyperVSites/%s/machines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HyperVSiteName, id.MachineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hyper V Site Machine ID
func (id HyperVSiteMachineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOffAzure", "Microsoft.OffAzure", "Microsoft.OffAzure"),
		resourceids.StaticSegment("staticHyperVSites", "hyperVSites", "hyperVSites"),
		resourceids.UserSpecifiedSegment("hyperVSiteName", "hyperVSiteValue"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineValue"),
	}
}

// String returns a human-readable description of this Hyper V Site Machine ID
func (id HyperVSiteMachineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hyper V Site Name: %q", id.HyperVSiteName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
	}
	return fmt.Sprintf("Hyper V Site Machine (%s)", strings.Join(components, "\n"))
}
