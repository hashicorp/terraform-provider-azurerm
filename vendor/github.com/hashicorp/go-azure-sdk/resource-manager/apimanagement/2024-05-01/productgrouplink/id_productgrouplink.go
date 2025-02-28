package productgrouplink

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProductGroupLinkId{})
}

var _ resourceids.ResourceId = &ProductGroupLinkId{}

// ProductGroupLinkId is a struct representing the Resource ID for a Product Group Link
type ProductGroupLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	ProductId         string
	GroupLinkId       string
}

// NewProductGroupLinkID returns a new ProductGroupLinkId struct
func NewProductGroupLinkID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, productId string, groupLinkId string) ProductGroupLinkId {
	return ProductGroupLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		ProductId:         productId,
		GroupLinkId:       groupLinkId,
	}
}

// ParseProductGroupLinkID parses 'input' into a ProductGroupLinkId
func ParseProductGroupLinkID(input string) (*ProductGroupLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductGroupLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductGroupLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProductGroupLinkIDInsensitively parses 'input' case-insensitively into a ProductGroupLinkId
// note: this method should only be used for API response data and not user input
func ParseProductGroupLinkIDInsensitively(input string) (*ProductGroupLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductGroupLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductGroupLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProductGroupLinkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.GroupLinkId, ok = input.Parsed["groupLinkId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupLinkId", input)
	}

	return nil
}

// ValidateProductGroupLinkID checks that 'input' can be parsed as a Product Group Link ID
func ValidateProductGroupLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProductGroupLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Product Group Link ID
func (id ProductGroupLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/products/%s/groupLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.ProductId, id.GroupLinkId)
}

// Segments returns a slice of Resource ID Segments which comprise this Product Group Link ID
func (id ProductGroupLinkId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticGroupLinks", "groupLinks", "groupLinks"),
		resourceids.UserSpecifiedSegment("groupLinkId", "groupLinkId"),
	}
}

// String returns a human-readable description of this Product Group Link ID
func (id ProductGroupLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Product: %q", id.ProductId),
		fmt.Sprintf("Group Link: %q", id.GroupLinkId),
	}
	return fmt.Sprintf("Product Group Link (%s)", strings.Join(components, "\n"))
}
