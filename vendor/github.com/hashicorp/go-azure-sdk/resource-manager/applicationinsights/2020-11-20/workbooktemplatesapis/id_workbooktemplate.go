package workbooktemplatesapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkbookTemplateId{})
}

var _ resourceids.ResourceId = &WorkbookTemplateId{}

// WorkbookTemplateId is a struct representing the Resource ID for a Workbook Template
type WorkbookTemplateId struct {
	SubscriptionId       string
	ResourceGroupName    string
	WorkbookTemplateName string
}

// NewWorkbookTemplateID returns a new WorkbookTemplateId struct
func NewWorkbookTemplateID(subscriptionId string, resourceGroupName string, workbookTemplateName string) WorkbookTemplateId {
	return WorkbookTemplateId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		WorkbookTemplateName: workbookTemplateName,
	}
}

// ParseWorkbookTemplateID parses 'input' into a WorkbookTemplateId
func ParseWorkbookTemplateID(input string) (*WorkbookTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkbookTemplateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkbookTemplateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkbookTemplateIDInsensitively parses 'input' case-insensitively into a WorkbookTemplateId
// note: this method should only be used for API response data and not user input
func ParseWorkbookTemplateIDInsensitively(input string) (*WorkbookTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkbookTemplateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkbookTemplateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkbookTemplateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkbookTemplateName, ok = input.Parsed["workbookTemplateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workbookTemplateName", input)
	}

	return nil
}

// ValidateWorkbookTemplateID checks that 'input' can be parsed as a Workbook Template ID
func ValidateWorkbookTemplateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkbookTemplateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workbook Template ID
func (id WorkbookTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/workbookTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkbookTemplateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Workbook Template ID
func (id WorkbookTemplateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticWorkbookTemplates", "workbookTemplates", "workbookTemplates"),
		resourceids.UserSpecifiedSegment("workbookTemplateName", "workbookTemplateName"),
	}
}

// String returns a human-readable description of this Workbook Template ID
func (id WorkbookTemplateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workbook Template Name: %q", id.WorkbookTemplateName),
	}
	return fmt.Sprintf("Workbook Template (%s)", strings.Join(components, "\n"))
}
