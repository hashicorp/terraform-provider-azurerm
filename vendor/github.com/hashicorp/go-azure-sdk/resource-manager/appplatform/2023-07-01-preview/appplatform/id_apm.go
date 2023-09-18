package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ApmId{}

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
	parser := resourceids.NewParserFromResourceIdType(ApmId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApmId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.ApmName, ok = parsed.Parsed["apmName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "apmName", *parsed)
	}

	return &id, nil
}

// ParseApmIDInsensitively parses 'input' case-insensitively into a ApmId
// note: this method should only be used for API response data and not user input
func ParseApmIDInsensitively(input string) (*ApmId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApmId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApmId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.ApmName, ok = parsed.Parsed["apmName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "apmName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticApms", "apms", "apms"),
		resourceids.UserSpecifiedSegment("apmName", "apmValue"),
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
