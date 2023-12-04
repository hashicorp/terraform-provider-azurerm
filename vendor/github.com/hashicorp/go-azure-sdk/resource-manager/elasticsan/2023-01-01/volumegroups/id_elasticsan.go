package volumegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ElasticSanId{}

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
	parser := resourceids.NewParserFromResourceIdType(ElasticSanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ElasticSanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ElasticSanName, ok = parsed.Parsed["elasticSanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "elasticSanName", *parsed)
	}

	return &id, nil
}

// ParseElasticSanIDInsensitively parses 'input' case-insensitively into a ElasticSanId
// note: this method should only be used for API response data and not user input
func ParseElasticSanIDInsensitively(input string) (*ElasticSanId, error) {
	parser := resourceids.NewParserFromResourceIdType(ElasticSanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ElasticSanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ElasticSanName, ok = parsed.Parsed["elasticSanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "elasticSanName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("elasticSanName", "elasticSanValue"),
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
