package configurationprofilehcrpassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&Providers2ConfigurationProfileAssignmentId{})
}

var _ resourceids.ResourceId = &Providers2ConfigurationProfileAssignmentId{}

// Providers2ConfigurationProfileAssignmentId is a struct representing the Resource ID for a Providers 2 Configuration Profile Assignment
type Providers2ConfigurationProfileAssignmentId struct {
	SubscriptionId                     string
	ResourceGroupName                  string
	MachineName                        string
	ConfigurationProfileAssignmentName string
}

// NewProviders2ConfigurationProfileAssignmentID returns a new Providers2ConfigurationProfileAssignmentId struct
func NewProviders2ConfigurationProfileAssignmentID(subscriptionId string, resourceGroupName string, machineName string, configurationProfileAssignmentName string) Providers2ConfigurationProfileAssignmentId {
	return Providers2ConfigurationProfileAssignmentId{
		SubscriptionId:                     subscriptionId,
		ResourceGroupName:                  resourceGroupName,
		MachineName:                        machineName,
		ConfigurationProfileAssignmentName: configurationProfileAssignmentName,
	}
}

// ParseProviders2ConfigurationProfileAssignmentID parses 'input' into a Providers2ConfigurationProfileAssignmentId
func ParseProviders2ConfigurationProfileAssignmentID(input string) (*Providers2ConfigurationProfileAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2ConfigurationProfileAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2ConfigurationProfileAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviders2ConfigurationProfileAssignmentIDInsensitively parses 'input' case-insensitively into a Providers2ConfigurationProfileAssignmentId
// note: this method should only be used for API response data and not user input
func ParseProviders2ConfigurationProfileAssignmentIDInsensitively(input string) (*Providers2ConfigurationProfileAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&Providers2ConfigurationProfileAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Providers2ConfigurationProfileAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Providers2ConfigurationProfileAssignmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ConfigurationProfileAssignmentName, ok = input.Parsed["configurationProfileAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationProfileAssignmentName", input)
	}

	return nil
}

// ValidateProviders2ConfigurationProfileAssignmentID checks that 'input' can be parsed as a Providers 2 Configuration Profile Assignment ID
func ValidateProviders2ConfigurationProfileAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviders2ConfigurationProfileAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Providers 2 Configuration Profile Assignment ID
func (id Providers2ConfigurationProfileAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/providers/Microsoft.AutoManage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MachineName, id.ConfigurationProfileAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Providers 2 Configuration Profile Assignment ID
func (id Providers2ConfigurationProfileAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHybridCompute", "Microsoft.HybridCompute", "Microsoft.HybridCompute"),
		resourceids.StaticSegment("staticMachines", "machines", "machines"),
		resourceids.UserSpecifiedSegment("machineName", "machineName"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutoManage", "Microsoft.AutoManage", "Microsoft.AutoManage"),
		resourceids.StaticSegment("staticConfigurationProfileAssignments", "configurationProfileAssignments", "configurationProfileAssignments"),
		resourceids.UserSpecifiedSegment("configurationProfileAssignmentName", "configurationProfileAssignmentName"),
	}
}

// String returns a human-readable description of this Providers 2 Configuration Profile Assignment ID
func (id Providers2ConfigurationProfileAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Machine Name: %q", id.MachineName),
		fmt.Sprintf("Configuration Profile Assignment Name: %q", id.ConfigurationProfileAssignmentName),
	}
	return fmt.Sprintf("Providers 2 Configuration Profile Assignment (%s)", strings.Join(components, "\n"))
}
