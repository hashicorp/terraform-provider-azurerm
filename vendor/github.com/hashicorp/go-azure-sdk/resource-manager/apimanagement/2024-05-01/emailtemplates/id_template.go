package emailtemplates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TemplateId{})
}

var _ resourceids.ResourceId = &TemplateId{}

// TemplateId is a struct representing the Resource ID for a Template
type TemplateId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	TemplateName      TemplateName
}

// NewTemplateID returns a new TemplateId struct
func NewTemplateID(subscriptionId string, resourceGroupName string, serviceName string, templateName TemplateName) TemplateId {
	return TemplateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		TemplateName:      templateName,
	}
}

// ParseTemplateID parses 'input' into a TemplateId
func ParseTemplateID(input string) (*TemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TemplateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TemplateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTemplateIDInsensitively parses 'input' case-insensitively into a TemplateId
// note: this method should only be used for API response data and not user input
func ParseTemplateIDInsensitively(input string) (*TemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TemplateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TemplateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TemplateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if v, ok := input.Parsed["templateName"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "templateName", input)
		}

		templateName, err := parseTemplateName(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.TemplateName = *templateName
	}

	return nil
}

// ValidateTemplateID checks that 'input' can be parsed as a Template ID
func ValidateTemplateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTemplateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Template ID
func (id TemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/templates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, string(id.TemplateName))
}

// Segments returns a slice of Resource ID Segments which comprise this Template ID
func (id TemplateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticTemplates", "templates", "templates"),
		resourceids.ConstantSegment("templateName", PossibleValuesForTemplateName(), "accountClosedDeveloper"),
	}
}

// String returns a human-readable description of this Template ID
func (id TemplateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Template Name: %q", string(id.TemplateName)),
	}
	return fmt.Sprintf("Template (%s)", strings.Join(components, "\n"))
}
