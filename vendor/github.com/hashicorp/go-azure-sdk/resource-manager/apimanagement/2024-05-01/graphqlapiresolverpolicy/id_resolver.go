package graphqlapiresolverpolicy

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResolverId{})
}

var _ resourceids.ResourceId = &ResolverId{}

// ResolverId is a struct representing the Resource ID for a Resolver
type ResolverId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ApiId             string
	ResolverId        string
}

// NewResolverID returns a new ResolverId struct
func NewResolverID(subscriptionId string, resourceGroupName string, serviceName string, apiId string, resolverId string) ResolverId {
	return ResolverId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ApiId:             apiId,
		ResolverId:        resolverId,
	}
}

// ParseResolverID parses 'input' into a ResolverId
func ParseResolverID(input string) (*ResolverId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResolverId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResolverId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResolverIDInsensitively parses 'input' case-insensitively into a ResolverId
// note: this method should only be used for API response data and not user input
func ParseResolverIDInsensitively(input string) (*ResolverId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResolverId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResolverId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResolverId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.ApiId, ok = input.Parsed["apiId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiId", input)
	}

	if id.ResolverId, ok = input.Parsed["resolverId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resolverId", input)
	}

	return nil
}

// ValidateResolverID checks that 'input' can be parsed as a Resolver ID
func ValidateResolverID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResolverID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resolver ID
func (id ResolverId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/resolvers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ApiId, id.ResolverId)
}

// Segments returns a slice of Resource ID Segments which comprise this Resolver ID
func (id ResolverId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiId", "apiId"),
		resourceids.StaticSegment("staticResolvers", "resolvers", "resolvers"),
		resourceids.UserSpecifiedSegment("resolverId", "resolverId"),
	}
}

// String returns a human-readable description of this Resolver ID
func (id ResolverId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Resolver: %q", id.ResolverId),
	}
	return fmt.Sprintf("Resolver (%s)", strings.Join(components, "\n"))
}
