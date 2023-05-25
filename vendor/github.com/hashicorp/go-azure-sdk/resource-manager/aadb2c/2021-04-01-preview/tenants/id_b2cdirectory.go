package tenants

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = B2CDirectoryId{}

// B2CDirectoryId is a struct representing the Resource ID for a B 2 C Directory
type B2CDirectoryId struct {
	SubscriptionId string
	ResourceGroup  string
	DirectoryName  string
}

// NewB2CDirectoryID returns a new B2CDirectoryId struct
func NewB2CDirectoryID(subscriptionId string, resourceGroup string, directoryName string) B2CDirectoryId {
	return B2CDirectoryId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DirectoryName:  directoryName,
	}
}

// ParseB2CDirectoryID parses 'input' into a B2CDirectoryId
func ParseB2CDirectoryID(input string) (*B2CDirectoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(B2CDirectoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := B2CDirectoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", *parsed)
	}

	if id.DirectoryName, ok = parsed.Parsed["directoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "directoryName", *parsed)
	}

	return &id, nil
}

// ParseB2CDirectoryIDInsensitively parses 'input' case-insensitively into a B2CDirectoryId
// note: this method should only be used for API response data and not user input
func ParseB2CDirectoryIDInsensitively(input string) (*B2CDirectoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(B2CDirectoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := B2CDirectoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", *parsed)
	}

	if id.DirectoryName, ok = parsed.Parsed["directoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "directoryName", *parsed)
	}

	return &id, nil
}

// ValidateB2CDirectoryID checks that 'input' can be parsed as a B 2 C Directory ID
func ValidateB2CDirectoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseB2CDirectoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted B 2 C Directory ID
func (id B2CDirectoryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureActiveDirectory/b2cDirectories/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DirectoryName)
}

// Segments returns a slice of Resource ID Segments which comprise this B 2 C Directory ID
func (id B2CDirectoryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("microsoftAzureActiveDirectory", "Microsoft.AzureActiveDirectory", "Microsoft.AzureActiveDirectory"),
		resourceids.StaticSegment("b2cDirectories", "b2cDirectories", "b2cDirectories"),
		resourceids.UserSpecifiedSegment("directoryName", "directoryValue"),
	}
}

// String returns a human-readable description of this B 2 C Directory ID
func (id B2CDirectoryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Directory Name: %q", id.DirectoryName),
	}
	return fmt.Sprintf("B 2 C Directory (%s)", strings.Join(components, "\n"))
}
