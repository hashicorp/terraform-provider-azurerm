package organizationresources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&OrganizationId{})
}

var _ resourceids.ResourceId = &OrganizationId{}

// OrganizationId is a struct representing the Resource ID for a Organization
type OrganizationId struct {
	SubscriptionId    string
	ResourceGroupName string
	OrganizationName  string
}

// NewOrganizationID returns a new OrganizationId struct
func NewOrganizationID(subscriptionId string, resourceGroupName string, organizationName string) OrganizationId {
	return OrganizationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		OrganizationName:  organizationName,
	}
}

// ParseOrganizationID parses 'input' into a OrganizationId
func ParseOrganizationID(input string) (*OrganizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OrganizationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OrganizationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseOrganizationIDInsensitively parses 'input' case-insensitively into a OrganizationId
// note: this method should only be used for API response data and not user input
func ParseOrganizationIDInsensitively(input string) (*OrganizationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&OrganizationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := OrganizationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *OrganizationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.OrganizationName, ok = input.Parsed["organizationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "organizationName", input)
	}

	return nil
}

// ValidateOrganizationID checks that 'input' can be parsed as a Organization ID
func ValidateOrganizationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOrganizationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Organization ID
func (id OrganizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Confluent/organizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.OrganizationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Organization ID
func (id OrganizationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftConfluent", "Microsoft.Confluent", "Microsoft.Confluent"),
		resourceids.StaticSegment("staticOrganizations", "organizations", "organizations"),
		resourceids.UserSpecifiedSegment("organizationName", "organizationName"),
	}
}

// String returns a human-readable description of this Organization ID
func (id OrganizationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Organization Name: %q", id.OrganizationName),
	}
	return fmt.Sprintf("Organization (%s)", strings.Join(components, "\n"))
}
