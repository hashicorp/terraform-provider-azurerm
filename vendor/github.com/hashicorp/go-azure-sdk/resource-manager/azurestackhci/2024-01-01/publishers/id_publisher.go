package publishers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PublisherId{})
}

var _ resourceids.ResourceId = &PublisherId{}

// PublisherId is a struct representing the Resource ID for a Publisher
type PublisherId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	PublisherName     string
}

// NewPublisherID returns a new PublisherId struct
func NewPublisherID(subscriptionId string, resourceGroupName string, clusterName string, publisherName string) PublisherId {
	return PublisherId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		PublisherName:     publisherName,
	}
}

// ParsePublisherID parses 'input' into a PublisherId
func ParsePublisherID(input string) (*PublisherId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PublisherId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PublisherId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePublisherIDInsensitively parses 'input' case-insensitively into a PublisherId
// note: this method should only be used for API response data and not user input
func ParsePublisherIDInsensitively(input string) (*PublisherId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PublisherId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PublisherId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PublisherId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	if id.PublisherName, ok = input.Parsed["publisherName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "publisherName", input)
	}

	return nil
}

// ValidatePublisherID checks that 'input' can be parsed as a Publisher ID
func ValidatePublisherID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePublisherID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Publisher ID
func (id PublisherId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/publishers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.PublisherName)
}

// Segments returns a slice of Resource ID Segments which comprise this Publisher ID
func (id PublisherId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticPublishers", "publishers", "publishers"),
		resourceids.UserSpecifiedSegment("publisherName", "publisherName"),
	}
}

// String returns a human-readable description of this Publisher ID
func (id PublisherId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Publisher Name: %q", id.PublisherName),
	}
	return fmt.Sprintf("Publisher (%s)", strings.Join(components, "\n"))
}
