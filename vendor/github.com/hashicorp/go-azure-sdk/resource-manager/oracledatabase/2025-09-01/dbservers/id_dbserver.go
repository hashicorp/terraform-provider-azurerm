package dbservers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DbServerId{})
}

var _ resourceids.ResourceId = &DbServerId{}

// DbServerId is a struct representing the Resource ID for a Db Server
type DbServerId struct {
	SubscriptionId                 string
	ResourceGroupName              string
	CloudExadataInfrastructureName string
	DbServerName                   string
}

// NewDbServerID returns a new DbServerId struct
func NewDbServerID(subscriptionId string, resourceGroupName string, cloudExadataInfrastructureName string, dbServerName string) DbServerId {
	return DbServerId{
		SubscriptionId:                 subscriptionId,
		ResourceGroupName:              resourceGroupName,
		CloudExadataInfrastructureName: cloudExadataInfrastructureName,
		DbServerName:                   dbServerName,
	}
}

// ParseDbServerID parses 'input' into a DbServerId
func ParseDbServerID(input string) (*DbServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDbServerIDInsensitively parses 'input' case-insensitively into a DbServerId
// note: this method should only be used for API response data and not user input
func ParseDbServerIDInsensitively(input string) (*DbServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbServerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DbServerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudExadataInfrastructureName, ok = input.Parsed["cloudExadataInfrastructureName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudExadataInfrastructureName", input)
	}

	if id.DbServerName, ok = input.Parsed["dbServerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dbServerName", input)
	}

	return nil
}

// ValidateDbServerID checks that 'input' can be parsed as a Db Server ID
func ValidateDbServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDbServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Db Server ID
func (id DbServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/cloudExadataInfrastructures/%s/dbServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudExadataInfrastructureName, id.DbServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Db Server ID
func (id DbServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticCloudExadataInfrastructures", "cloudExadataInfrastructures", "cloudExadataInfrastructures"),
		resourceids.UserSpecifiedSegment("cloudExadataInfrastructureName", "cloudExadataInfrastructureName"),
		resourceids.StaticSegment("staticDbServers", "dbServers", "dbServers"),
		resourceids.UserSpecifiedSegment("dbServerName", "dbServerName"),
	}
}

// String returns a human-readable description of this Db Server ID
func (id DbServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Exadata Infrastructure Name: %q", id.CloudExadataInfrastructureName),
		fmt.Sprintf("Db Server Name: %q", id.DbServerName),
	}
	return fmt.Sprintf("Db Server (%s)", strings.Join(components, "\n"))
}
