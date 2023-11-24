package objectdatatypes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ObjectDataTypeId{}

// ObjectDataTypeId is a struct representing the Resource ID for a Object Data Type
type ObjectDataTypeId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	ObjectDataTypeName    string
}

// NewObjectDataTypeID returns a new ObjectDataTypeId struct
func NewObjectDataTypeID(subscriptionId string, resourceGroupName string, automationAccountName string, objectDataTypeName string) ObjectDataTypeId {
	return ObjectDataTypeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		ObjectDataTypeName:    objectDataTypeName,
	}
}

// ParseObjectDataTypeID parses 'input' into a ObjectDataTypeId
func ParseObjectDataTypeID(input string) (*ObjectDataTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ObjectDataTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ObjectDataTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.ObjectDataTypeName, ok = parsed.Parsed["objectDataTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "objectDataTypeName", *parsed)
	}

	return &id, nil
}

// ParseObjectDataTypeIDInsensitively parses 'input' case-insensitively into a ObjectDataTypeId
// note: this method should only be used for API response data and not user input
func ParseObjectDataTypeIDInsensitively(input string) (*ObjectDataTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(ObjectDataTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ObjectDataTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.ObjectDataTypeName, ok = parsed.Parsed["objectDataTypeName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "objectDataTypeName", *parsed)
	}

	return &id, nil
}

// ValidateObjectDataTypeID checks that 'input' can be parsed as a Object Data Type ID
func ValidateObjectDataTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseObjectDataTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Object Data Type ID
func (id ObjectDataTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/objectDataTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.ObjectDataTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Object Data Type ID
func (id ObjectDataTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticObjectDataTypes", "objectDataTypes", "objectDataTypes"),
		resourceids.UserSpecifiedSegment("objectDataTypeName", "objectDataTypeValue"),
	}
}

// String returns a human-readable description of this Object Data Type ID
func (id ObjectDataTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Object Data Type Name: %q", id.ObjectDataTypeName),
	}
	return fmt.Sprintf("Object Data Type (%s)", strings.Join(components, "\n"))
}
