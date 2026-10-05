package sapdatabaseinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabaseInstanceId{})
}

var _ resourceids.ResourceId = &DatabaseInstanceId{}

// DatabaseInstanceId is a struct representing the Resource ID for a Database Instance
type DatabaseInstanceId struct {
	SubscriptionId         string
	ResourceGroupName      string
	SapVirtualInstanceName string
	DatabaseInstanceName   string
}

// NewDatabaseInstanceID returns a new DatabaseInstanceId struct
func NewDatabaseInstanceID(subscriptionId string, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string) DatabaseInstanceId {
	return DatabaseInstanceId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		SapVirtualInstanceName: sapVirtualInstanceName,
		DatabaseInstanceName:   databaseInstanceName,
	}
}

// ParseDatabaseInstanceID parses 'input' into a DatabaseInstanceId
func ParseDatabaseInstanceID(input string) (*DatabaseInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseInstanceIDInsensitively parses 'input' case-insensitively into a DatabaseInstanceId
// note: this method should only be used for API response data and not user input
func ParseDatabaseInstanceIDInsensitively(input string) (*DatabaseInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SapVirtualInstanceName, ok = input.Parsed["sapVirtualInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sapVirtualInstanceName", input)
	}

	if id.DatabaseInstanceName, ok = input.Parsed["databaseInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseInstanceName", input)
	}

	return nil
}

// ValidateDatabaseInstanceID checks that 'input' can be parsed as a Database Instance ID
func ValidateDatabaseInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Instance ID
func (id DatabaseInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Workloads/sapVirtualInstances/%s/databaseInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName, id.DatabaseInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Instance ID
func (id DatabaseInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWorkloads", "Microsoft.Workloads", "Microsoft.Workloads"),
		resourceids.StaticSegment("staticSapVirtualInstances", "sapVirtualInstances", "sapVirtualInstances"),
		resourceids.UserSpecifiedSegment("sapVirtualInstanceName", "sapVirtualInstanceName"),
		resourceids.StaticSegment("staticDatabaseInstances", "databaseInstances", "databaseInstances"),
		resourceids.UserSpecifiedSegment("databaseInstanceName", "databaseInstanceName"),
	}
}

// String returns a human-readable description of this Database Instance ID
func (id DatabaseInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sap Virtual Instance Name: %q", id.SapVirtualInstanceName),
		fmt.Sprintf("Database Instance Name: %q", id.DatabaseInstanceName),
	}
	return fmt.Sprintf("Database Instance (%s)", strings.Join(components, "\n"))
}
