package backend

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BackendId{})
}

var _ resourceids.ResourceId = &BackendId{}

// BackendId is a struct representing the Resource ID for a Backend
type BackendId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	BackendId         string
}

// NewBackendID returns a new BackendId struct
func NewBackendID(subscriptionId string, resourceGroupName string, serviceName string, backendId string) BackendId {
	return BackendId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		BackendId:         backendId,
	}
}

// ParseBackendID parses 'input' into a BackendId
func ParseBackendID(input string) (*BackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackendId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackendIDInsensitively parses 'input' case-insensitively into a BackendId
// note: this method should only be used for API response data and not user input
func ParseBackendIDInsensitively(input string) (*BackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackendId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackendId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BackendId, ok = input.Parsed["backendId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backendId", input)
	}

	return nil
}

// ValidateBackendID checks that 'input' can be parsed as a Backend ID
func ValidateBackendID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackendID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backend ID
func (id BackendId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/backends/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.BackendId)
}

// Segments returns a slice of Resource ID Segments which comprise this Backend ID
func (id BackendId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticBackends", "backends", "backends"),
		resourceids.UserSpecifiedSegment("backendId", "backendId"),
	}
}

// String returns a human-readable description of this Backend ID
func (id BackendId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Backend: %q", id.BackendId),
	}
	return fmt.Sprintf("Backend (%s)", strings.Join(components, "\n"))
}
