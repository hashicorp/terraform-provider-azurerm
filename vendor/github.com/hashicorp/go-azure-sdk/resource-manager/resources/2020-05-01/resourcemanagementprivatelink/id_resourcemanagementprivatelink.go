package resourcemanagementprivatelink

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResourceManagementPrivateLinkId{})
}

var _ resourceids.ResourceId = &ResourceManagementPrivateLinkId{}

// ResourceManagementPrivateLinkId is a struct representing the Resource ID for a Resource Management Private Link
type ResourceManagementPrivateLinkId struct {
	SubscriptionId                    string
	ResourceGroupName                 string
	ResourceManagementPrivateLinkName string
}

// NewResourceManagementPrivateLinkID returns a new ResourceManagementPrivateLinkId struct
func NewResourceManagementPrivateLinkID(subscriptionId string, resourceGroupName string, resourceManagementPrivateLinkName string) ResourceManagementPrivateLinkId {
	return ResourceManagementPrivateLinkId{
		SubscriptionId:                    subscriptionId,
		ResourceGroupName:                 resourceGroupName,
		ResourceManagementPrivateLinkName: resourceManagementPrivateLinkName,
	}
}

// ParseResourceManagementPrivateLinkID parses 'input' into a ResourceManagementPrivateLinkId
func ParseResourceManagementPrivateLinkID(input string) (*ResourceManagementPrivateLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceManagementPrivateLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceManagementPrivateLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResourceManagementPrivateLinkIDInsensitively parses 'input' case-insensitively into a ResourceManagementPrivateLinkId
// note: this method should only be used for API response data and not user input
func ParseResourceManagementPrivateLinkIDInsensitively(input string) (*ResourceManagementPrivateLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceManagementPrivateLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceManagementPrivateLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResourceManagementPrivateLinkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ResourceManagementPrivateLinkName, ok = input.Parsed["resourceManagementPrivateLinkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceManagementPrivateLinkName", input)
	}

	return nil
}

// ValidateResourceManagementPrivateLinkID checks that 'input' can be parsed as a Resource Management Private Link ID
func ValidateResourceManagementPrivateLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceManagementPrivateLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Management Private Link ID
func (id ResourceManagementPrivateLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Authorization/resourceManagementPrivateLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceManagementPrivateLinkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Management Private Link ID
func (id ResourceManagementPrivateLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticResourceManagementPrivateLinks", "resourceManagementPrivateLinks", "resourceManagementPrivateLinks"),
		resourceids.UserSpecifiedSegment("resourceManagementPrivateLinkName", "resourceManagementPrivateLinkName"),
	}
}

// String returns a human-readable description of this Resource Management Private Link ID
func (id ResourceManagementPrivateLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Management Private Link Name: %q", id.ResourceManagementPrivateLinkName),
	}
	return fmt.Sprintf("Resource Management Private Link (%s)", strings.Join(components, "\n"))
}
