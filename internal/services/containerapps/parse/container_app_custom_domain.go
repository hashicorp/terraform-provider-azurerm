// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &ContainerAppCustomDomainId{}

// ContainerAppCustomDomainId is a struct representing the Resource ID for a Container App
type ContainerAppCustomDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	ContainerAppName  string
	CustomDomainName  string
}

// NewContainerAppCustomDomainId returns a new ContainerAppCustomDomainId struct
func NewContainerAppCustomDomainId(subscriptionId string, resourceGroupName string, containerAppName string, customDomainName string) ContainerAppCustomDomainId {
	return ContainerAppCustomDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ContainerAppName:  containerAppName,
		CustomDomainName:  customDomainName,
	}
}

// ContainerAppCustomDomainID parses 'input' into a ContainerAppCustomDomainId
func ContainerAppCustomDomainID(input string) (*ContainerAppCustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerAppCustomDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerAppCustomDomainId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ContainerAppCustomDomainIDInsensitively parses 'input' case-insensitively into a ContainerAppCustomDomainId
// note: this method should only be used for API response data and not user input
func ContainerAppCustomDomainIDInsensitively(input string) (*ContainerAppCustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ContainerAppCustomDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ContainerAppCustomDomainId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ContainerAppCustomDomainId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ContainerAppName, ok = input.Parsed["containerAppName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerAppName", input)
	}

	if id.CustomDomainName, ok = input.Parsed["customDomainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "customDomainName", input)
	}

	return nil
}

// ID returns the formatted Container App ID
func (id ContainerAppCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/containerApps/%s/customDomainName/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName, id.CustomDomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container App ID
func (id ContainerAppCustomDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticContainerApps", "containerApps", "containerApps"),
		resourceids.UserSpecifiedSegment("containerAppName", "containerAppValue"),
		resourceids.StaticSegment("staticCustomDomainName", "customDomainName", "customDomainName"),
		resourceids.UserSpecifiedSegment("customDomainName", "customDomainNameValue"),
	}
}

// String returns a human-readable description of this Container App ID
func (id ContainerAppCustomDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container App Name: %q", id.ContainerAppName),
		fmt.Sprintf("Custom Domain Name: %q", id.CustomDomainName),
	}
	return fmt.Sprintf("Container App Custom Domain (%s)", strings.Join(components, "\n"))
}
