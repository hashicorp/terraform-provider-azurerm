package objectdatatypes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ModuleObjectDataTypeId{})
}

var _ resourceids.ResourceId = &ModuleObjectDataTypeId{}

// ModuleObjectDataTypeId is a struct representing the Resource ID for a Module Object Data Type
type ModuleObjectDataTypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	ModuleName            string
	ObjectDataTypeName    string
}

// NewModuleObjectDataTypeID returns a new ModuleObjectDataTypeId struct
func NewModuleObjectDataTypeID(subscriptionId string, resourceGroupName string, automationAccountName string, moduleName string, objectDataTypeName string) ModuleObjectDataTypeId {
	return ModuleObjectDataTypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		ModuleName:            moduleName,
		ObjectDataTypeName:    objectDataTypeName,
	}
}

// ParseModuleObjectDataTypeID parses 'input' into a ModuleObjectDataTypeId
func ParseModuleObjectDataTypeID(input string) (*ModuleObjectDataTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ModuleObjectDataTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ModuleObjectDataTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseModuleObjectDataTypeIDInsensitively parses 'input' case-insensitively into a ModuleObjectDataTypeId
// note: this method should only be used for API response data and not user input
func ParseModuleObjectDataTypeIDInsensitively(input string) (*ModuleObjectDataTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ModuleObjectDataTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ModuleObjectDataTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ModuleObjectDataTypeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutomationAccountName, ok = input.Parsed["automationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", input)
	}

	if id.ModuleName, ok = input.Parsed["moduleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "moduleName", input)
	}

	if id.ObjectDataTypeName, ok = input.Parsed["objectDataTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "objectDataTypeName", input)
	}

	return nil
}

// ValidateModuleObjectDataTypeID checks that 'input' can be parsed as a Module Object Data Type ID
func ValidateModuleObjectDataTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseModuleObjectDataTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Module Object Data Type ID
func (id ModuleObjectDataTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/modules/%s/objectDataTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.ModuleName, id.ObjectDataTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Module Object Data Type ID
func (id ModuleObjectDataTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticModules", "modules", "modules"),
		resourceids.UserSpecifiedSegment("moduleName", "moduleName"),
		resourceids.StaticSegment("staticObjectDataTypes", "objectDataTypes", "objectDataTypes"),
		resourceids.UserSpecifiedSegment("objectDataTypeName", "objectDataTypeName"),
	}
}

// String returns a human-readable description of this Module Object Data Type ID
func (id ModuleObjectDataTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Module Name: %q", id.ModuleName),
		fmt.Sprintf("Object Data Type Name: %q", id.ObjectDataTypeName),
	}
	return fmt.Sprintf("Module Object Data Type (%s)", strings.Join(components, "\n"))
}
