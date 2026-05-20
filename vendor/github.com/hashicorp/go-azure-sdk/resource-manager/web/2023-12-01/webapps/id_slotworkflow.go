package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SlotWorkflowId{})
}

var _ resourceids.ResourceId = &SlotWorkflowId{}

// SlotWorkflowId is a struct representing the Resource ID for a Slot Workflow
type SlotWorkflowId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	WorkflowName      string
}

// NewSlotWorkflowID returns a new SlotWorkflowId struct
func NewSlotWorkflowID(subscriptionId string, resourceGroupName string, siteName string, slotName string, workflowName string) SlotWorkflowId {
	return SlotWorkflowId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		WorkflowName:      workflowName,
	}
}

// ParseSlotWorkflowID parses 'input' into a SlotWorkflowId
func ParseSlotWorkflowID(input string) (*SlotWorkflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotWorkflowId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotWorkflowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotWorkflowIDInsensitively parses 'input' case-insensitively into a SlotWorkflowId
// note: this method should only be used for API response data and not user input
func ParseSlotWorkflowIDInsensitively(input string) (*SlotWorkflowId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotWorkflowId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotWorkflowId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotWorkflowId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.WorkflowName, ok = input.Parsed["workflowName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workflowName", input)
	}

	return nil
}

// ValidateSlotWorkflowID checks that 'input' can be parsed as a Slot Workflow ID
func ValidateSlotWorkflowID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotWorkflowID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Workflow ID
func (id SlotWorkflowId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/workflows/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.WorkflowName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Workflow ID
func (id SlotWorkflowId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotName"),
		resourceids.StaticSegment("staticWorkflows", "workflows", "workflows"),
		resourceids.UserSpecifiedSegment("workflowName", "workflowName"),
	}
}

// String returns a human-readable description of this Slot Workflow ID
func (id SlotWorkflowId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Workflow Name: %q", id.WorkflowName),
	}
	return fmt.Sprintf("Slot Workflow (%s)", strings.Join(components, "\n"))
}
