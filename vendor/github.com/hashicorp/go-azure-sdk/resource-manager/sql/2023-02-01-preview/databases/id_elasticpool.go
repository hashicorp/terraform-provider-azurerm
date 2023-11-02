package databases

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ElasticPoolId{}

// ElasticPoolId is a struct representing the Resource ID for a Elastic Pool
type ElasticPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	ElasticPoolName   string
}

// NewElasticPoolID returns a new ElasticPoolId struct
func NewElasticPoolID(subscriptionId string, resourceGroupName string, serverName string, elasticPoolName string) ElasticPoolId {
	return ElasticPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		ElasticPoolName:   elasticPoolName,
	}
}

// ParseElasticPoolID parses 'input' into a ElasticPoolId
func ParseElasticPoolID(input string) (*ElasticPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(ElasticPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ElasticPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverName", *parsed)
	}

	if id.ElasticPoolName, ok = parsed.Parsed["elasticPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "elasticPoolName", *parsed)
	}

	return &id, nil
}

// ParseElasticPoolIDInsensitively parses 'input' case-insensitively into a ElasticPoolId
// note: this method should only be used for API response data and not user input
func ParseElasticPoolIDInsensitively(input string) (*ElasticPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(ElasticPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ElasticPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serverName", *parsed)
	}

	if id.ElasticPoolName, ok = parsed.Parsed["elasticPoolName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "elasticPoolName", *parsed)
	}

	return &id, nil
}

// ValidateElasticPoolID checks that 'input' can be parsed as a Elastic Pool ID
func ValidateElasticPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseElasticPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Elastic Pool ID
func (id ElasticPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/elasticPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.ElasticPoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Elastic Pool ID
func (id ElasticPoolId) Segments() []resourceids.Segment {
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

// String returns a human-readable description of this Elastic Pool ID
func (id ElasticPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Elastic Pool Name: %q", id.ElasticPoolName),
	}
	return fmt.Sprintf("Elastic Pool (%s)", strings.Join(components, "\n"))
}
