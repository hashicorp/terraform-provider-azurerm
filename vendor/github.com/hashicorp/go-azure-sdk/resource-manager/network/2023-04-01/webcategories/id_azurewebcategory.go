package webcategories

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AzureWebCategoryId{}

// AzureWebCategoryId is a struct representing the Resource ID for a Azure Web Category
type AzureWebCategoryId struct {
	SubscriptionId       string
	AzureWebCategoryName string
}

// NewAzureWebCategoryID returns a new AzureWebCategoryId struct
func NewAzureWebCategoryID(subscriptionId string, azureWebCategoryName string) AzureWebCategoryId {
	return AzureWebCategoryId{
		SubscriptionId:       subscriptionId,
		AzureWebCategoryName: azureWebCategoryName,
	}
}

// ParseAzureWebCategoryID parses 'input' into a AzureWebCategoryId
func ParseAzureWebCategoryID(input string) (*AzureWebCategoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(AzureWebCategoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AzureWebCategoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.AzureWebCategoryName, ok = parsed.Parsed["azureWebCategoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "azureWebCategoryName", *parsed)
	}

	return &id, nil
}

// ParseAzureWebCategoryIDInsensitively parses 'input' case-insensitively into a AzureWebCategoryId
// note: this method should only be used for API response data and not user input
func ParseAzureWebCategoryIDInsensitively(input string) (*AzureWebCategoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(AzureWebCategoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AzureWebCategoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.AzureWebCategoryName, ok = parsed.Parsed["azureWebCategoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "azureWebCategoryName", *parsed)
	}

	return &id, nil
}

// ValidateAzureWebCategoryID checks that 'input' can be parsed as a Azure Web Category ID
func ValidateAzureWebCategoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAzureWebCategoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Azure Web Category ID
func (id AzureWebCategoryId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/azureWebCategories/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AzureWebCategoryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Azure Web Category ID
func (id AzureWebCategoryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticAzureWebCategories", "azureWebCategories", "azureWebCategories"),
		resourceids.UserSpecifiedSegment("azureWebCategoryName", "azureWebCategoryValue"),
	}
}

// String returns a human-readable description of this Azure Web Category ID
func (id AzureWebCategoryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Azure Web Category Name: %q", id.AzureWebCategoryName),
	}
	return fmt.Sprintf("Azure Web Category (%s)", strings.Join(components, "\n"))
}
