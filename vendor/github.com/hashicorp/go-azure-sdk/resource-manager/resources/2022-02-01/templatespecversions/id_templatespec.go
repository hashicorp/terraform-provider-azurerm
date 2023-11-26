package templatespecversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TemplateSpecId{}

// TemplateSpecId is a struct representing the Resource ID for a Template Spec
type TemplateSpecId struct {
	SubscriptionId    string
	ResourceGroupName string
	TemplateSpecName  string
}

// NewTemplateSpecID returns a new TemplateSpecId struct
func NewTemplateSpecID(subscriptionId string, resourceGroupName string, templateSpecName string) TemplateSpecId {
	return TemplateSpecId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		TemplateSpecName:  templateSpecName,
	}
}

// ParseTemplateSpecID parses 'input' into a TemplateSpecId
func ParseTemplateSpecID(input string) (*TemplateSpecId, error) {
	parser := resourceids.NewParserFromResourceIdType(TemplateSpecId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TemplateSpecId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TemplateSpecName, ok = parsed.Parsed["templateSpecName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "templateSpecName", *parsed)
	}

	return &id, nil
}

// ParseTemplateSpecIDInsensitively parses 'input' case-insensitively into a TemplateSpecId
// note: this method should only be used for API response data and not user input
func ParseTemplateSpecIDInsensitively(input string) (*TemplateSpecId, error) {
	parser := resourceids.NewParserFromResourceIdType(TemplateSpecId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TemplateSpecId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TemplateSpecName, ok = parsed.Parsed["templateSpecName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "templateSpecName", *parsed)
	}

	return &id, nil
}

// ValidateTemplateSpecID checks that 'input' can be parsed as a Template Spec ID
func ValidateTemplateSpecID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTemplateSpecID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Template Spec ID
func (id TemplateSpecId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/templateSpecs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TemplateSpecName)
}

// Segments returns a slice of Resource ID Segments which comprise this Template Spec ID
func (id TemplateSpecId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticTemplateSpecs", "templateSpecs", "templateSpecs"),
		resourceids.UserSpecifiedSegment("templateSpecName", "templateSpecValue"),
	}
}

// String returns a human-readable description of this Template Spec ID
func (id TemplateSpecId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Template Spec Name: %q", id.TemplateSpecName),
	}
	return fmt.Sprintf("Template Spec (%s)", strings.Join(components, "\n"))
}
