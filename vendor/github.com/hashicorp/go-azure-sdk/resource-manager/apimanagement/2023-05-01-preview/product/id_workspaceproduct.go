package product

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &WorkspaceProductId{}

// WorkspaceProductId is a struct representing the Resource ID for a Workspace Product
type WorkspaceProductId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	ProductId         string
}

// NewWorkspaceProductID returns a new WorkspaceProductId struct
func NewWorkspaceProductID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, productId string) WorkspaceProductId {
	return WorkspaceProductId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		ProductId:         productId,
	}
}

// ParseWorkspaceProductID parses 'input' into a WorkspaceProductId
func ParseWorkspaceProductID(input string) (*WorkspaceProductId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceProductId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceProductId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceProductIDInsensitively parses 'input' case-insensitively into a WorkspaceProductId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceProductIDInsensitively(input string) (*WorkspaceProductId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceProductId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceProductId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceProductId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateWorkspaceProductID checks that 'input' can be parsed as a Workspace Product ID
func ValidateWorkspaceProductID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceProductID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Product ID
func (id WorkspaceProductId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/products/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.ProductId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Product ID
func (id WorkspaceProductId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceIdValue"),
		resourceids.StaticSegment("staticProducts", "products", "products"),
		resourceids.UserSpecifiedSegment("productId", "productIdValue"),
	}
}

// String returns a human-readable description of this Workspace Product ID
func (id WorkspaceProductId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Product: %q", id.ProductId),
	}
	return fmt.Sprintf("Workspace Product (%s)", strings.Join(components, "\n"))
}
