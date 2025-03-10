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
	recaser.RegisterResourceId(&ProductApiLinkId{})
}

var _ resourceids.ResourceId = &ProductApiLinkId{}

// ProductApiLinkId is a struct representing the Resource ID for a Product Api Link
type ProductApiLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ProductId         string
	ApiLinkId         string
}

// NewProductApiLinkID returns a new ProductApiLinkId struct
func NewProductApiLinkID(subscriptionId string, resourceGroupName string, serviceName string, productId string, apiLinkId string) ProductApiLinkId {
	return ProductApiLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ProductId:         productId,
		ApiLinkId:         apiLinkId,
	}
}

// ParseProductApiLinkID parses 'input' into a ProductApiLinkId
func ParseProductApiLinkID(input string) (*ProductApiLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductApiLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductApiLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProductApiLinkIDInsensitively parses 'input' case-insensitively into a ProductApiLinkId
// note: this method should only be used for API response data and not user input
func ParseProductApiLinkIDInsensitively(input string) (*ProductApiLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductApiLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductApiLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProductApiLinkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ProductId, ok = input.Parsed["productId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "productId", input)
	}

	if id.ApiLinkId, ok = input.Parsed["apiLinkId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiLinkId", input)
	}

	return nil
}

// ValidateProductApiLinkID checks that 'input' can be parsed as a Product Api Link ID
func ValidateProductApiLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProductApiLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Product Api Link ID
func (id ProductApiLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/apiLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ProductId, id.ApiLinkId)
}

// Segments returns a slice of Resource ID Segments which comprise this Product Api Link ID
func (id ProductApiLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticProducts", "products", "products"),
		resourceids.UserSpecifiedSegment("productId", "productId"),
		resourceids.StaticSegment("staticApiLinks", "apiLinks", "apiLinks"),
		resourceids.UserSpecifiedSegment("apiLinkId", "apiLinkId"),
	}
}

// String returns a human-readable description of this Product Api Link ID
func (id ProductApiLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Product: %q", id.ProductId),
		fmt.Sprintf("Api Link: %q", id.ApiLinkId),
	}
	return fmt.Sprintf("Product Api Link (%s)", strings.Join(components, "\n"))
}
