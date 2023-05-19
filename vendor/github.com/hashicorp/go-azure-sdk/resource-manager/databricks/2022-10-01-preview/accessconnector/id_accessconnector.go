package accessconnector

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AccessConnectorId{}

// AccessConnectorId is a struct representing the Resource ID for a Access Connector
type AccessConnectorId struct {
	SubscriptionId      string
	ResourceGroupName   string
	AccessConnectorName string
}

// NewAccessConnectorID returns a new AccessConnectorId struct
func NewAccessConnectorID(subscriptionId string, resourceGroupName string, accessConnectorName string) AccessConnectorId {
	return AccessConnectorId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		AccessConnectorName: accessConnectorName,
	}
}

// ParseAccessConnectorID parses 'input' into a AccessConnectorId
func ParseAccessConnectorID(input string) (*AccessConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessConnectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessConnectorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccessConnectorName, ok = parsed.Parsed["accessConnectorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accessConnectorName", *parsed)
	}

	return &id, nil
}

// ParseAccessConnectorIDInsensitively parses 'input' case-insensitively into a AccessConnectorId
// note: this method should only be used for API response data and not user input
func ParseAccessConnectorIDInsensitively(input string) (*AccessConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(AccessConnectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AccessConnectorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccessConnectorName, ok = parsed.Parsed["accessConnectorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accessConnectorName", *parsed)
	}

	return &id, nil
}

// ValidateAccessConnectorID checks that 'input' can be parsed as a Access Connector ID
func ValidateAccessConnectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccessConnectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Access Connector ID
func (id AccessConnectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Databricks/accessConnectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccessConnectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Access Connector ID
func (id AccessConnectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDatabricks", "Microsoft.Databricks", "Microsoft.Databricks"),
		resourceids.StaticSegment("staticAccessConnectors", "accessConnectors", "accessConnectors"),
		resourceids.UserSpecifiedSegment("accessConnectorName", "accessConnectorValue"),
	}
}

// String returns a human-readable description of this Access Connector ID
func (id AccessConnectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Access Connector Name: %q", id.AccessConnectorName),
	}
	return fmt.Sprintf("Access Connector (%s)", strings.Join(components, "\n"))
}
