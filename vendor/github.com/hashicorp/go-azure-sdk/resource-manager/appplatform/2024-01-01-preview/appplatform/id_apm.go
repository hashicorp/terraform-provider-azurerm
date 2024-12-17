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
	recaser.RegisterResourceId(&ApmId{})
}

var _ resourceids.ResourceId = &ApmId{}

// ApmId is a struct representing the Resource ID for a Apm
type ApmId struct {
	SubscriptionId    string
	ResourceGroupName string
	SpringName        string
	ApmName           string
}

// NewApmID returns a new ApmId struct
func NewApmID(subscriptionId string, resourceGroupName string, springName string, apmName string) ApmId {
	return ApmId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SpringName:        springName,
		ApmName:           apmName,
	}
}

// ParseApmID parses 'input' into a ApmId
func ParseApmID(input string) (*ApmId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApmId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApmId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseApmIDInsensitively parses 'input' case-insensitively into a ApmId
// note: this method should only be used for API response data and not user input
func ParseApmIDInsensitively(input string) (*ApmId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ApmId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ApmId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ApmId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ApmName, ok = input.Parsed["apmName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apmName", input)
	}

	return nil
}

// ValidateApmID checks that 'input' can be parsed as a Apm ID
func ValidateApmID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApmID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Apm ID
func (id ApmId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/apms/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ApmName)
}

// Segments returns a slice of Resource ID Segments which comprise this Apm ID
func (id ApmId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticApms", "apms", "apms"),
		resourceids.UserSpecifiedSegment("apmName", "apmName"),
	}
}

// String returns a human-readable description of this Apm ID
func (id ApmId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Apm Name: %q", id.ApmName),
	}
	return fmt.Sprintf("Apm (%s)", strings.Join(components, "\n"))
}
