package productapi

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProductApiId{})
}

var _ resourceids.ResourceId = &ProductApiId{}

// ProductApiId is a struct representing the Resource ID for a Product Api
type ProductApiId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ProductId         string
	ApiId             string
}

// NewProductApiID returns a new ProductApiId struct
func NewProductApiID(subscriptionId string, resourceGroupName string, serviceName string, productId string, apiId string) ProductApiId {
	return ProductApiId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ProductId:         productId,
		ApiId:             apiId,
	}
}

// ParseProductApiID parses 'input' into a ProductApiId
func ParseProductApiID(input string) (*ProductApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductApiId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProductApiIDInsensitively parses 'input' case-insensitively into a ProductApiId
// note: this method should only be used for API response data and not user input
func ParseProductApiIDInsensitively(input string) (*ProductApiId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductApiId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductApiId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProductApiId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ApiId, ok = input.Parsed["apiId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiId", input)
	}

	return nil
}

// ValidateProductApiID checks that 'input' can be parsed as a Product Api ID
func ValidateProductApiID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProductApiID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Product Api ID
func (id ProductApiId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/apis/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ProductId, id.ApiId)
}

// Segments returns a slice of Resource ID Segments which comprise this Product Api ID
func (id ProductApiId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiId", "apiId"),
	}
}

// String returns a human-readable description of this Product Api ID
func (id ProductApiId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Product: %q", id.ProductId),
		fmt.Sprintf("Api: %q", id.ApiId),
	}
	return fmt.Sprintf("Product Api (%s)", strings.Join(components, "\n"))
}
