package templatespecversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TemplateSpecVersionId{})
}

var _ resourceids.ResourceId = &TemplateSpecVersionId{}

// TemplateSpecVersionId is a struct representing the Resource ID for a Template Spec Version
type TemplateSpecVersionId struct {
	SubscriptionId    string
	ResourceGroupName string
	TemplateSpecName  string
	VersionName       string
}

// NewTemplateSpecVersionID returns a new TemplateSpecVersionId struct
func NewTemplateSpecVersionID(subscriptionId string, resourceGroupName string, templateSpecName string, versionName string) TemplateSpecVersionId {
	return TemplateSpecVersionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		TemplateSpecName:  templateSpecName,
		VersionName:       versionName,
	}
}

// ParseTemplateSpecVersionID parses 'input' into a TemplateSpecVersionId
func ParseTemplateSpecVersionID(input string) (*TemplateSpecVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TemplateSpecVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TemplateSpecVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTemplateSpecVersionIDInsensitively parses 'input' case-insensitively into a TemplateSpecVersionId
// note: this method should only be used for API response data and not user input
func ParseTemplateSpecVersionIDInsensitively(input string) (*TemplateSpecVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TemplateSpecVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TemplateSpecVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TemplateSpecVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.TemplateSpecName, ok = input.Parsed["templateSpecName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "templateSpecName", input)
	}

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	return nil
}

// ValidateTemplateSpecVersionID checks that 'input' can be parsed as a Template Spec Version ID
func ValidateTemplateSpecVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTemplateSpecVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Template Spec Version ID
func (id TemplateSpecVersionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/templateSpecs/%s/versions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TemplateSpecName, id.VersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Template Spec Version ID
func (id TemplateSpecVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticTemplateSpecs", "templateSpecs", "templateSpecs"),
		resourceids.UserSpecifiedSegment("templateSpecName", "templateSpecName"),
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionName"),
	}
}

// String returns a human-readable description of this Template Spec Version ID
func (id TemplateSpecVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Template Spec Name: %q", id.TemplateSpecName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
	}
	return fmt.Sprintf("Template Spec Version (%s)", strings.Join(components, "\n"))
}
