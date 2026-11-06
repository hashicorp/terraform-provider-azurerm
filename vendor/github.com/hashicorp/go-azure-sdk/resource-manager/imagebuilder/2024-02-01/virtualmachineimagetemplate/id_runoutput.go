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
	recaser.RegisterResourceId(&RunOutputId{})
}

var _ resourceids.ResourceId = &RunOutputId{}

// RunOutputId is a struct representing the Resource ID for a Run Output
type RunOutputId struct {
	SubscriptionId    string
	ResourceGroupName string
	ImageTemplateName string
	RunOutputName     string
}

// NewRunOutputID returns a new RunOutputId struct
func NewRunOutputID(subscriptionId string, resourceGroupName string, imageTemplateName string, runOutputName string) RunOutputId {
	return RunOutputId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ImageTemplateName: imageTemplateName,
		RunOutputName:     runOutputName,
	}
}

// ParseRunOutputID parses 'input' into a RunOutputId
func ParseRunOutputID(input string) (*RunOutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RunOutputId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RunOutputId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRunOutputIDInsensitively parses 'input' case-insensitively into a RunOutputId
// note: this method should only be used for API response data and not user input
func ParseRunOutputIDInsensitively(input string) (*RunOutputId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RunOutputId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RunOutputId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RunOutputId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RunOutputName, ok = input.Parsed["runOutputName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runOutputName", input)
	}

	return nil
}

// ValidateRunOutputID checks that 'input' can be parsed as a Run Output ID
func ValidateRunOutputID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRunOutputID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Run Output ID
func (id RunOutputId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.VirtualMachineImages/imageTemplates/%s/runOutputs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ImageTemplateName, id.RunOutputName)
}

// Segments returns a slice of Resource ID Segments which comprise this Run Output ID
func (id RunOutputId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftVirtualMachineImages", "Microsoft.VirtualMachineImages", "Microsoft.VirtualMachineImages"),
		resourceids.StaticSegment("staticImageTemplates", "imageTemplates", "imageTemplates"),
		resourceids.UserSpecifiedSegment("imageTemplateName", "imageTemplateName"),
		resourceids.StaticSegment("staticRunOutputs", "runOutputs", "runOutputs"),
		resourceids.UserSpecifiedSegment("runOutputName", "runOutputName"),
	}
}

// String returns a human-readable description of this Run Output ID
func (id RunOutputId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Image Template Name: %q", id.ImageTemplateName),
		fmt.Sprintf("Run Output Name: %q", id.RunOutputName),
	}
	return fmt.Sprintf("Run Output (%s)", strings.Join(components, "\n"))
}
