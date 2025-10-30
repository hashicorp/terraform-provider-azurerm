package typefields

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TypeId{})
}

var _ resourceids.ResourceId = &TypeId{}

// TypeId is a struct representing the Resource ID for a Type
type TypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	ModuleName            string
	TypeName              string
}

// NewTypeID returns a new TypeId struct
func NewTypeID(subscriptionId string, resourceGroupName string, automationAccountName string, moduleName string, typeName string) TypeId {
	return TypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		ModuleName:            moduleName,
		TypeName:              typeName,
	}
}

// ParseTypeID parses 'input' into a TypeId
func ParseTypeID(input string) (*TypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTypeIDInsensitively parses 'input' case-insensitively into a TypeId
// note: this method should only be used for API response data and not user input
func ParseTypeIDInsensitively(input string) (*TypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TypeName, ok = input.Parsed["typeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "typeName", input)
	}

	return nil
}

// ValidateTypeID checks that 'input' can be parsed as a Type ID
func ValidateTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Type ID
func (id TypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/modules/%s/types/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.ModuleName, id.TypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Type ID
func (id TypeId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticTypes", "types", "types"),
		resourceids.UserSpecifiedSegment("typeName", "typeName"),
	}
}

// String returns a human-readable description of this Type ID
func (id TypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Module Name: %q", id.ModuleName),
		fmt.Sprintf("Type Name: %q", id.TypeName),
	}
	return fmt.Sprintf("Type (%s)", strings.Join(components, "\n"))
}
