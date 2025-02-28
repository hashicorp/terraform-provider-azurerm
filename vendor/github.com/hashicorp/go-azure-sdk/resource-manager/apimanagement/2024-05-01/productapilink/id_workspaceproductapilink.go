package productapilink

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceProductApiLinkId{})
}

var _ resourceids.ResourceId = &WorkspaceProductApiLinkId{}

// WorkspaceProductApiLinkId is a struct representing the Resource ID for a Workspace Product Api Link
type WorkspaceProductApiLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	ProductId         string
	ApiLinkId         string
}

// NewWorkspaceProductApiLinkID returns a new WorkspaceProductApiLinkId struct
func NewWorkspaceProductApiLinkID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, productId string, apiLinkId string) WorkspaceProductApiLinkId {
	return WorkspaceProductApiLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		ProductId:         productId,
		ApiLinkId:         apiLinkId,
	}
}

// ParseWorkspaceProductApiLinkID parses 'input' into a WorkspaceProductApiLinkId
func ParseWorkspaceProductApiLinkID(input string) (*WorkspaceProductApiLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceProductApiLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceProductApiLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceProductApiLinkIDInsensitively parses 'input' case-insensitively into a WorkspaceProductApiLinkId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceProductApiLinkIDInsensitively(input string) (*WorkspaceProductApiLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceProductApiLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceProductApiLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceProductApiLinkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.WorkspaceId, ok = input.Parsed["workspaceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceId", input)
	}

	if id.ProductId, ok = input.Parsed["productId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "productId", input)
	}

	if id.ApiLinkId, ok = input.Parsed["apiLinkId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiLinkId", input)
	}

	return nil
}

// ValidateWorkspaceProductApiLinkID checks that 'input' can be parsed as a Workspace Product Api Link ID
func ValidateWorkspaceProductApiLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceProductApiLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Product Api Link ID
func (id WorkspaceProductApiLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/products/%s/apiLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.ProductId, id.ApiLinkId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Product Api Link ID
func (id WorkspaceProductApiLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceId"),
		resourceids.StaticSegment("staticProducts", "products", "products"),
		resourceids.UserSpecifiedSegment("productId", "productId"),
		resourceids.StaticSegment("staticApiLinks", "apiLinks", "apiLinks"),
		resourceids.UserSpecifiedSegment("apiLinkId", "apiLinkId"),
	}
}

// String returns a human-readable description of this Workspace Product Api Link ID
func (id WorkspaceProductApiLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Product: %q", id.ProductId),
		fmt.Sprintf("Api Link: %q", id.ApiLinkId),
	}
	return fmt.Sprintf("Workspace Product Api Link (%s)", strings.Join(components, "\n"))
}
