package virtualmachineimagetemplate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ImageTemplateId{})
}

var _ resourceids.ResourceId = &ImageTemplateId{}

// ImageTemplateId is a struct representing the Resource ID for a Image Template
type ImageTemplateId struct {
	SubscriptionId    string
	ResourceGroupName string
	ImageTemplateName string
}

// NewImageTemplateID returns a new ImageTemplateId struct
func NewImageTemplateID(subscriptionId string, resourceGroupName string, imageTemplateName string) ImageTemplateId {
	return ImageTemplateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ImageTemplateName: imageTemplateName,
	}
}

// ParseImageTemplateID parses 'input' into a ImageTemplateId
func ParseImageTemplateID(input string) (*ImageTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageTemplateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageTemplateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseImageTemplateIDInsensitively parses 'input' case-insensitively into a ImageTemplateId
// note: this method should only be used for API response data and not user input
func ParseImageTemplateIDInsensitively(input string) (*ImageTemplateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ImageTemplateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ImageTemplateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ImageTemplateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ImageTemplateName, ok = input.Parsed["imageTemplateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "imageTemplateName", input)
	}

	return nil
}

// ValidateImageTemplateID checks that 'input' can be parsed as a Image Template ID
func ValidateImageTemplateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseImageTemplateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Image Template ID
func (id ImageTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.VirtualMachineImages/imageTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ImageTemplateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Image Template ID
func (id ImageTemplateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftVirtualMachineImages", "Microsoft.VirtualMachineImages", "Microsoft.VirtualMachineImages"),
		resourceids.StaticSegment("staticImageTemplates", "imageTemplates", "imageTemplates"),
		resourceids.UserSpecifiedSegment("imageTemplateName", "imageTemplateName"),
	}
}

// String returns a human-readable description of this Image Template ID
func (id ImageTemplateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Image Template Name: %q", id.ImageTemplateName),
	}
	return fmt.Sprintf("Image Template (%s)", strings.Join(components, "\n"))
}
