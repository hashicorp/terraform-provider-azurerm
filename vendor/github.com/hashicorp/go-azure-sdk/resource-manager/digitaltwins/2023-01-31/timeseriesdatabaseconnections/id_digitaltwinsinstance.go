package timeseriesdatabaseconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DigitalTwinsInstanceId{})
}

var _ resourceids.ResourceId = &DigitalTwinsInstanceId{}

// DigitalTwinsInstanceId is a struct representing the Resource ID for a Digital Twins Instance
type DigitalTwinsInstanceId struct {
	SubscriptionId           string
	ResourceGroupName        string
	DigitalTwinsInstanceName string
}

// NewDigitalTwinsInstanceID returns a new DigitalTwinsInstanceId struct
func NewDigitalTwinsInstanceID(subscriptionId string, resourceGroupName string, digitalTwinsInstanceName string) DigitalTwinsInstanceId {
	return DigitalTwinsInstanceId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		DigitalTwinsInstanceName: digitalTwinsInstanceName,
	}
}

// ParseDigitalTwinsInstanceID parses 'input' into a DigitalTwinsInstanceId
func ParseDigitalTwinsInstanceID(input string) (*DigitalTwinsInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DigitalTwinsInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DigitalTwinsInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDigitalTwinsInstanceIDInsensitively parses 'input' case-insensitively into a DigitalTwinsInstanceId
// note: this method should only be used for API response data and not user input
func ParseDigitalTwinsInstanceIDInsensitively(input string) (*DigitalTwinsInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DigitalTwinsInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DigitalTwinsInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DigitalTwinsInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DigitalTwinsInstanceName, ok = input.Parsed["digitalTwinsInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "digitalTwinsInstanceName", input)
	}

	return nil
}

// ValidateDigitalTwinsInstanceID checks that 'input' can be parsed as a Digital Twins Instance ID
func ValidateDigitalTwinsInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDigitalTwinsInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Digital Twins Instance ID
func (id DigitalTwinsInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DigitalTwinsInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Digital Twins Instance ID
func (id DigitalTwinsInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDigitalTwins", "Microsoft.DigitalTwins", "Microsoft.DigitalTwins"),
		resourceids.StaticSegment("staticDigitalTwinsInstances", "digitalTwinsInstances", "digitalTwinsInstances"),
		resourceids.UserSpecifiedSegment("digitalTwinsInstanceName", "digitalTwinsInstanceName"),
	}
}

// String returns a human-readable description of this Digital Twins Instance ID
func (id DigitalTwinsInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Digital Twins Instance Name: %q", id.DigitalTwinsInstanceName),
	}
	return fmt.Sprintf("Digital Twins Instance (%s)", strings.Join(components, "\n"))
}
