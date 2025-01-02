package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CustomDomainId{})
}

var _ resourceids.ResourceId = &CustomDomainId{}

// CustomDomainId is a struct representing the Resource ID for a Custom Domain
type CustomDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	SignalRName       string
	CustomDomainName  string
}

// NewCustomDomainID returns a new CustomDomainId struct
func NewCustomDomainID(subscriptionId string, resourceGroupName string, signalRName string, customDomainName string) CustomDomainId {
	return CustomDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SignalRName:       signalRName,
		CustomDomainName:  customDomainName,
	}
}

// ParseCustomDomainID parses 'input' into a CustomDomainId
func ParseCustomDomainID(input string) (*CustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CustomDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CustomDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCustomDomainIDInsensitively parses 'input' case-insensitively into a CustomDomainId
// note: this method should only be used for API response data and not user input
func ParseCustomDomainIDInsensitively(input string) (*CustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CustomDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CustomDomainId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CustomDomainId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SignalRName, ok = input.Parsed["signalRName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "signalRName", input)
	}

	if id.CustomDomainName, ok = input.Parsed["customDomainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "customDomainName", input)
	}

	return nil
}

// ValidateCustomDomainID checks that 'input' can be parsed as a Custom Domain ID
func ValidateCustomDomainID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Custom Domain ID
func (id CustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/signalR/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SignalRName, id.CustomDomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Custom Domain ID
func (id CustomDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticSignalR", "signalR", "signalR"),
		resourceids.UserSpecifiedSegment("signalRName", "signalRName"),
		resourceids.StaticSegment("staticCustomDomains", "customDomains", "customDomains"),
		resourceids.UserSpecifiedSegment("customDomainName", "customDomainName"),
	}
}

// String returns a human-readable description of this Custom Domain ID
func (id CustomDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Signal R Name: %q", id.SignalRName),
		fmt.Sprintf("Custom Domain Name: %q", id.CustomDomainName),
	}
	return fmt.Sprintf("Custom Domain (%s)", strings.Join(components, "\n"))
}
