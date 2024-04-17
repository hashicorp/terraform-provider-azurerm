package encodings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TransformId{})
}

var _ resourceids.ResourceId = &TransformId{}

// TransformId is a struct representing the Resource ID for a Transform
type TransformId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	TransformName     string
}

// NewTransformID returns a new TransformId struct
func NewTransformID(subscriptionId string, resourceGroupName string, mediaServiceName string, transformName string) TransformId {
	return TransformId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		TransformName:     transformName,
	}
}

// ParseTransformID parses 'input' into a TransformId
func ParseTransformID(input string) (*TransformId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TransformId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TransformId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTransformIDInsensitively parses 'input' case-insensitively into a TransformId
// note: this method should only be used for API response data and not user input
func ParseTransformIDInsensitively(input string) (*TransformId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TransformId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TransformId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TransformId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MediaServiceName, ok = input.Parsed["mediaServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", input)
	}

	if id.TransformName, ok = input.Parsed["transformName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "transformName", input)
	}

	return nil
}

// ValidateTransformID checks that 'input' can be parsed as a Transform ID
func ValidateTransformID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTransformID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Transform ID
func (id TransformId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/transforms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.TransformName)
}

// Segments returns a slice of Resource ID Segments which comprise this Transform ID
func (id TransformId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticTransforms", "transforms", "transforms"),
		resourceids.UserSpecifiedSegment("transformName", "transformValue"),
	}
}

// String returns a human-readable description of this Transform ID
func (id TransformId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Transform Name: %q", id.TransformName),
	}
	return fmt.Sprintf("Transform (%s)", strings.Join(components, "\n"))
}
