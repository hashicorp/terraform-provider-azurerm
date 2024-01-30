// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SqlElasticPoolId{}

// SqlElasticPoolId is a struct representing the Resource ID for a Sql SqlElastic Pool
type SqlElasticPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	ElasticPoolName   string
}

// NewSqlElasticPoolID returns a new SqlElasticPoolId struct
func NewSqlElasticPoolID(subscriptionId string, resourceGroupName string, serverName string, elasticPoolName string) SqlElasticPoolId {
	return SqlElasticPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		ElasticPoolName:   elasticPoolName,
	}
}

// ParseSqlElasticPoolID parses 'input' into a SqlElasticPoolId
func ParseSqlElasticPoolID(input string) (*SqlElasticPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlElasticPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlElasticPoolId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSqlElasticPoolIDInsensitively parses 'input' case-insensitively into a SqlElasticPoolId
// note: this method should only be used for API response data and not user input
func ParseSqlElasticPoolIDInsensitively(input string) (*SqlElasticPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SqlElasticPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SqlElasticPoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SqlElasticPoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerName, ok = input.Parsed["serverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverName", input)
	}

	if id.ElasticPoolName, ok = input.Parsed["elasticPoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "elasticPoolName", input)
	}

	return nil
}

// ValidateSqlElasticPoolID checks that 'input' can be parsed as a Sql Elastic Pool ID
func ValidateSqlElasticPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlElasticPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Elastic Pool ID
func (id SqlElasticPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/elasticPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.ElasticPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Elastic Pool ID
func (id SqlElasticPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverValue"),
		resourceids.StaticSegment("staticElasticPools", "elasticPools", "elasticPools"),
		resourceids.UserSpecifiedSegment("elasticPoolName", "elasticPoolValue"),
	}
}

// String returns a human-readable description of this Sql Elastic Pool ID
func (id SqlElasticPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Elastic Pool Name: %q", id.ElasticPoolName),
	}
	return fmt.Sprintf("Elastic Pool (%s)", strings.Join(components, "\n"))
}
