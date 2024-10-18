package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ApplicationAcceleratorId{})
}

var _ resourceids.ResourceId = &ApplicationAcceleratorId{}

// ApplicationAcceleratorId is a struct representing the Resource ID for a Application Accelerator
type ApplicationAcceleratorId struct {
	SubscriptionId             string
	ResourceGroupName          string
	SpringName                 string
	ApplicationAcceleratorName string
}

// NewApplicationAcceleratorID returns a new ApplicationAcceleratorId struct
func NewApplicationAcceleratorID(subscriptionId string, resourceGroupName string, springName string, applicationAcceleratorName string) ApplicationAcceleratorId {
	return ApplicationAcceleratorId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		SpringName:                 springName,
		ApplicationAcceleratorName: applicationAcceleratorName,
	}
}

// ParseApplicationAcceleratorID parses 'input' into a ApplicationAcceleratorId
func ParseApplicationAcceleratorID(input string) (*ApplicationAcceleratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationAcceleratorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationAcceleratorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApplicationAcceleratorIDInsensitively parses 'input' case-insensitively into a ApplicationAcceleratorId
// note: this method should only be used for API response data and not user input
func ParseApplicationAcceleratorIDInsensitively(input string) (*ApplicationAcceleratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApplicationAcceleratorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApplicationAcceleratorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApplicationAcceleratorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.ApplicationAcceleratorName, ok = input.Parsed["applicationAcceleratorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "applicationAcceleratorName", input)
	}

	return nil
}

// ValidateApplicationAcceleratorID checks that 'input' can be parsed as a Application Accelerator ID
func ValidateApplicationAcceleratorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationAcceleratorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Accelerator ID
func (id ApplicationAcceleratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/applicationAccelerators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApplicationAcceleratorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Accelerator ID
func (id ApplicationAcceleratorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApplicationAccelerators", "applicationAccelerators", "applicationAccelerators"),
		resourceids.UserSpecifiedSegment("applicationAcceleratorName", "applicationAcceleratorName"),
	}
}

// String returns a human-readable description of this Application Accelerator ID
func (id ApplicationAcceleratorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Application Accelerator Name: %q", id.ApplicationAcceleratorName),
	}
	return fmt.Sprintf("Application Accelerator (%s)", strings.Join(components, "\n"))
}
