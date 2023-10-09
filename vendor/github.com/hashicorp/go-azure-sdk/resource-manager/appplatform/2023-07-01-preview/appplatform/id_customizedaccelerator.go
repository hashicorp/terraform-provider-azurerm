package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CustomizedAcceleratorId{}

// CustomizedAcceleratorId is a struct representing the Resource ID for a Customized Accelerator
type CustomizedAcceleratorId struct {
	SubscriptionId             string
	ResourceGroupName          string
	SpringName                 string
	ApplicationAcceleratorName string
	CustomizedAcceleratorName  string
}

// NewCustomizedAcceleratorID returns a new CustomizedAcceleratorId struct
func NewCustomizedAcceleratorID(subscriptionId string, resourceGroupName string, springName string, applicationAcceleratorName string, customizedAcceleratorName string) CustomizedAcceleratorId {
	return CustomizedAcceleratorId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		SpringName:                 springName,
		ApplicationAcceleratorName: applicationAcceleratorName,
		CustomizedAcceleratorName:  customizedAcceleratorName,
	}
}

// ParseCustomizedAcceleratorID parses 'input' into a CustomizedAcceleratorId
func ParseCustomizedAcceleratorID(input string) (*CustomizedAcceleratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomizedAcceleratorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomizedAcceleratorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.ApplicationAcceleratorName, ok = parsed.Parsed["applicationAcceleratorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applicationAcceleratorName", *parsed)
	}

	if id.CustomizedAcceleratorName, ok = parsed.Parsed["customizedAcceleratorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "customizedAcceleratorName", *parsed)
	}

	return &id, nil
}

// ParseCustomizedAcceleratorIDInsensitively parses 'input' case-insensitively into a CustomizedAcceleratorId
// note: this method should only be used for API response data and not user input
func ParseCustomizedAcceleratorIDInsensitively(input string) (*CustomizedAcceleratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomizedAcceleratorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomizedAcceleratorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.ApplicationAcceleratorName, ok = parsed.Parsed["applicationAcceleratorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "applicationAcceleratorName", *parsed)
	}

	if id.CustomizedAcceleratorName, ok = parsed.Parsed["customizedAcceleratorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "customizedAcceleratorName", *parsed)
	}

	return &id, nil
}

// ValidateCustomizedAcceleratorID checks that 'input' can be parsed as a Customized Accelerator ID
func ValidateCustomizedAcceleratorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomizedAcceleratorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Customized Accelerator ID
func (id CustomizedAcceleratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/applicationAccelerators/%s/customizedAccelerators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Customized Accelerator ID
func (id CustomizedAcceleratorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticApplicationAccelerators", "applicationAccelerators", "applicationAccelerators"),
		resourceids.UserSpecifiedSegment("applicationAcceleratorName", "applicationAcceleratorValue"),
		resourceids.StaticSegment("staticCustomizedAccelerators", "customizedAccelerators", "customizedAccelerators"),
		resourceids.UserSpecifiedSegment("customizedAcceleratorName", "customizedAcceleratorValue"),
	}
}

// String returns a human-readable description of this Customized Accelerator ID
func (id CustomizedAcceleratorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Application Accelerator Name: %q", id.ApplicationAcceleratorName),
		fmt.Sprintf("Customized Accelerator Name: %q", id.CustomizedAcceleratorName),
	}
	return fmt.Sprintf("Customized Accelerator (%s)", strings.Join(components, "\n"))
}
