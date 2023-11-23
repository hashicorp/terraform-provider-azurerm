package virtualmachinetemplates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualMachineTemplateId{}

// VirtualMachineTemplateId is a struct representing the Resource ID for a Virtual Machine Template
type VirtualMachineTemplateId struct {
	SubscriptionId             string
	ResourceGroupName          string
	VirtualMachineTemplateName string
}

// NewVirtualMachineTemplateID returns a new VirtualMachineTemplateId struct
func NewVirtualMachineTemplateID(subscriptionId string, resourceGroupName string, virtualMachineTemplateName string) VirtualMachineTemplateId {
	return VirtualMachineTemplateId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		VirtualMachineTemplateName: virtualMachineTemplateName,
	}
}

// ParseVirtualMachineTemplateID parses 'input' into a VirtualMachineTemplateId
func ParseVirtualMachineTemplateID(input string) (*VirtualMachineTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualMachineTemplateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualMachineTemplateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualMachineTemplateName, ok = parsed.Parsed["virtualMachineTemplateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineTemplateName", *parsed)
	}

	return &id, nil
}

// ParseVirtualMachineTemplateIDInsensitively parses 'input' case-insensitively into a VirtualMachineTemplateId
// note: this method should only be used for API response data and not user input
func ParseVirtualMachineTemplateIDInsensitively(input string) (*VirtualMachineTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualMachineTemplateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualMachineTemplateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VirtualMachineTemplateName, ok = parsed.Parsed["virtualMachineTemplateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualMachineTemplateName", *parsed)
	}

	return &id, nil
}

// ValidateVirtualMachineTemplateID checks that 'input' can be parsed as a Virtual Machine Template ID
func ValidateVirtualMachineTemplateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualMachineTemplateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Machine Template ID
func (id VirtualMachineTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ScVmm/virtualMachineTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualMachineTemplateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Machine Template ID
func (id VirtualMachineTemplateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftScVmm", "Microsoft.ScVmm", "Microsoft.ScVmm"),
		resourceids.StaticSegment("staticVirtualMachineTemplates", "virtualMachineTemplates", "virtualMachineTemplates"),
		resourceids.UserSpecifiedSegment("virtualMachineTemplateName", "virtualMachineTemplateValue"),
	}
}

// String returns a human-readable description of this Virtual Machine Template ID
func (id VirtualMachineTemplateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Machine Template Name: %q", id.VirtualMachineTemplateName),
	}
	return fmt.Sprintf("Virtual Machine Template (%s)", strings.Join(components, "\n"))
}
