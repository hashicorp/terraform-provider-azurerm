package users

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MongoClusterId{})
}

var _ resourceids.ResourceId = &MongoClusterId{}

// MongoClusterId is a struct representing the Resource ID for a Mongo Cluster
type MongoClusterId struct {
	SubscriptionId    string
	ResourceGroupName string
	MongoClusterName  string
}

// NewMongoClusterID returns a new MongoClusterId struct
func NewMongoClusterID(subscriptionId string, resourceGroupName string, mongoClusterName string) MongoClusterId {
	return MongoClusterId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MongoClusterName:  mongoClusterName,
	}
}

// ParseMongoClusterID parses 'input' into a MongoClusterId
func ParseMongoClusterID(input string) (*MongoClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongoClusterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongoClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMongoClusterIDInsensitively parses 'input' case-insensitively into a MongoClusterId
// note: this method should only be used for API response data and not user input
func ParseMongoClusterIDInsensitively(input string) (*MongoClusterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongoClusterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongoClusterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MongoClusterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MongoClusterName, ok = input.Parsed["mongoClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mongoClusterName", input)
	}

	return nil
}

// ValidateMongoClusterID checks that 'input' can be parsed as a Mongo Cluster ID
func ValidateMongoClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongoClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongo Cluster ID
func (id MongoClusterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/mongoClusters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MongoClusterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongo Cluster ID
func (id MongoClusterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticMongoClusters", "mongoClusters", "mongoClusters"),
		resourceids.UserSpecifiedSegment("mongoClusterName", "mongoClusterName"),
	}
}

// String returns a human-readable description of this Mongo Cluster ID
func (id MongoClusterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Mongo Cluster Name: %q", id.MongoClusterName),
	}
	return fmt.Sprintf("Mongo Cluster (%s)", strings.Join(components, "\n"))
}
