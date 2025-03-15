package producttag

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProductTagId{})
}

var _ resourceids.ResourceId = &ProductTagId{}

// ProductTagId is a struct representing the Resource ID for a Product Tag
type ProductTagId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ProductId         string
	TagId             string
}

// NewProductTagID returns a new ProductTagId struct
func NewProductTagID(subscriptionId string, resourceGroupName string, serviceName string, productId string, tagId string) ProductTagId {
	return ProductTagId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ProductId:         productId,
		TagId:             tagId,
	}
}

// ParseProductTagID parses 'input' into a ProductTagId
func ParseProductTagID(input string) (*ProductTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductTagId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProductTagIDInsensitively parses 'input' case-insensitively into a ProductTagId
// note: this method should only be used for API response data and not user input
func ParseProductTagIDInsensitively(input string) (*ProductTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProductTagId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProductTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProductTagId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TagId, ok = input.Parsed["tagId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagId", input)
	}

	return nil
}

// ValidateProductTagID checks that 'input' can be parsed as a Product Tag ID
func ValidateProductTagID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProductTagID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Product Tag ID
func (id ProductTagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/tags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ProductId, id.TagId)
}

// Segments returns a slice of Resource ID Segments which comprise this Product Tag ID
func (id ProductTagId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticTags", "tags", "tags"),
		resourceids.UserSpecifiedSegment("tagId", "tagId"),
	}
}

// String returns a human-readable description of this Product Tag ID
func (id ProductTagId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Product: %q", id.ProductId),
		fmt.Sprintf("Tag: %q", id.TagId),
	}
	return fmt.Sprintf("Product Tag (%s)", strings.Join(components, "\n"))
}
