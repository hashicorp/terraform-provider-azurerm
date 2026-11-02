package connectors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectorId{})
}

var _ resourceids.ResourceId = &ConnectorId{}

// ConnectorId is a struct representing the Resource ID for a Connector
type ConnectorId struct {
	SubscriptionId     string
	ResourceGroupName  string
	StorageAccountName string
	ConnectorName      string
}

// NewConnectorID returns a new ConnectorId struct
func NewConnectorID(subscriptionId string, resourceGroupName string, storageAccountName string, connectorName string) ConnectorId {
	return ConnectorId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		StorageAccountName: storageAccountName,
		ConnectorName:      connectorName,
	}
}

// ParseConnectorID parses 'input' into a ConnectorId
func ParseConnectorID(input string) (*ConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectorIDInsensitively parses 'input' case-insensitively into a ConnectorId
// note: this method should only be used for API response data and not user input
func ParseConnectorIDInsensitively(input string) (*ConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageAccountName, ok = input.Parsed["storageAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageAccountName", input)
	}

	if id.ConnectorName, ok = input.Parsed["connectorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectorName", input)
	}

	return nil
}

// ValidateConnectorID checks that 'input' can be parsed as a Connector ID
func ValidateConnectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connector ID
func (id ConnectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/connectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName, id.ConnectorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connector ID
func (id ConnectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorage", "Microsoft.Storage", "Microsoft.Storage"),
		resourceids.StaticSegment("staticStorageAccounts", "storageAccounts", "storageAccounts"),
		resourceids.UserSpecifiedSegment("storageAccountName", "storageAccountName"),
		resourceids.StaticSegment("staticConnectors", "connectors", "connectors"),
		resourceids.UserSpecifiedSegment("connectorName", "connectorName"),
	}
}

// String returns a human-readable description of this Connector ID
func (id ConnectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Account Name: %q", id.StorageAccountName),
		fmt.Sprintf("Connector Name: %q", id.ConnectorName),
	}
	return fmt.Sprintf("Connector (%s)", strings.Join(components, "\n"))
}
