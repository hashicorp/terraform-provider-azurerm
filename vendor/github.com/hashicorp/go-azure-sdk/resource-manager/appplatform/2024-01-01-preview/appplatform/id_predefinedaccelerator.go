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
	recaser.RegisterResourceId(&PredefinedAcceleratorId{})
}

var _ resourceids.ResourceId = &PredefinedAcceleratorId{}

// PredefinedAcceleratorId is a struct representing the Resource ID for a Predefined Accelerator
type PredefinedAcceleratorId struct {
	SubscriptionId             string
	ResourceGroupName          string
	SpringName                 string
	ApplicationAcceleratorName string
	PredefinedAcceleratorName  string
}

// NewPredefinedAcceleratorID returns a new PredefinedAcceleratorId struct
func NewPredefinedAcceleratorID(subscriptionId string, resourceGroupName string, springName string, applicationAcceleratorName string, predefinedAcceleratorName string) PredefinedAcceleratorId {
	return PredefinedAcceleratorId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		SpringName:                 springName,
		ApplicationAcceleratorName: applicationAcceleratorName,
		PredefinedAcceleratorName:  predefinedAcceleratorName,
	}
}

// ParsePredefinedAcceleratorID parses 'input' into a PredefinedAcceleratorId
func ParsePredefinedAcceleratorID(input string) (*PredefinedAcceleratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PredefinedAcceleratorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PredefinedAcceleratorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePredefinedAcceleratorIDInsensitively parses 'input' case-insensitively into a PredefinedAcceleratorId
// note: this method should only be used for API response data and not user input
func ParsePredefinedAcceleratorIDInsensitively(input string) (*PredefinedAcceleratorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PredefinedAcceleratorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PredefinedAcceleratorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PredefinedAcceleratorId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PredefinedAcceleratorName, ok = input.Parsed["predefinedAcceleratorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "predefinedAcceleratorName", input)
	}

	return nil
}

// ValidatePredefinedAcceleratorID checks that 'input' can be parsed as a Predefined Accelerator ID
func ValidatePredefinedAcceleratorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePredefinedAcceleratorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Predefined Accelerator ID
func (id PredefinedAcceleratorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/applicationAccelerators/%s/predefinedAccelerators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApplicationAcceleratorName, id.PredefinedAcceleratorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Predefined Accelerator ID
func (id PredefinedAcceleratorId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticPredefinedAccelerators", "predefinedAccelerators", "predefinedAccelerators"),
		resourceids.UserSpecifiedSegment("predefinedAcceleratorName", "predefinedAcceleratorName"),
	}
}

// String returns a human-readable description of this Predefined Accelerator ID
func (id PredefinedAcceleratorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Application Accelerator Name: %q", id.ApplicationAcceleratorName),
		fmt.Sprintf("Predefined Accelerator Name: %q", id.PredefinedAcceleratorName),
	}
	return fmt.Sprintf("Predefined Accelerator (%s)", strings.Join(components, "\n"))
}
