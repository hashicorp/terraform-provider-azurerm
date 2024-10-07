package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ElasticSanId{})
}

var _ resourceids.ResourceId = &ElasticSanId{}

// ElasticSanId is a struct representing the Resource ID for a Elastic San
type ElasticSanId struct {
	SubscriptionId    string
	ResourceGroupName string
	ElasticSanName    string
}

// NewElasticSanID returns a new ElasticSanId struct
func NewElasticSanID(subscriptionId string, resourceGroupName string, elasticSanName string) ElasticSanId {
	return ElasticSanId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ElasticSanName:    elasticSanName,
	}
}

// ParseElasticSanID parses 'input' into a ElasticSanId
func ParseElasticSanID(input string) (*ElasticSanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ElasticSanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ElasticSanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseElasticSanIDInsensitively parses 'input' case-insensitively into a ElasticSanId
// note: this method should only be used for API response data and not user input
func ParseElasticSanIDInsensitively(input string) (*ElasticSanId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ElasticSanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ElasticSanId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ElasticSanId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ElasticSanName, ok = input.Parsed["elasticSanName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "elasticSanName", input)
	}

	return nil
}

// ValidateElasticSanID checks that 'input' can be parsed as a Elastic San ID
func ValidateElasticSanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseElasticSanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Elastic San ID
func (id ElasticSanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ElasticSan/elasticSans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ElasticSanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Elastic San ID
func (id ElasticSanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftElasticSan", "Microsoft.ElasticSan", "Microsoft.ElasticSan"),
		resourceids.StaticSegment("staticElasticSans", "elasticSans", "elasticSans"),
		resourceids.UserSpecifiedSegment("elasticSanName", "elasticSanName"),
	}
}

// String returns a human-readable description of this Elastic San ID
func (id ElasticSanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Elastic San Name: %q", id.ElasticSanName),
	}
	return fmt.Sprintf("Elastic San (%s)", strings.Join(components, "\n"))
}
