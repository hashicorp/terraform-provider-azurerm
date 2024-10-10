package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&InstanceId{})
}

var _ resourceids.ResourceId = &InstanceId{}

// InstanceId is a struct representing the Resource ID for a Instance
type InstanceId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	InstanceId        string
}

// NewInstanceID returns a new InstanceId struct
func NewInstanceID(subscriptionId string, resourceGroupName string, siteName string, instanceId string) InstanceId {
	return InstanceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		InstanceId:        instanceId,
	}
}

// ParseInstanceID parses 'input' into a InstanceId
func ParseInstanceID(input string) (*InstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseInstanceIDInsensitively parses 'input' case-insensitively into a InstanceId
// note: this method should only be used for API response data and not user input
func ParseInstanceIDInsensitively(input string) (*InstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&InstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := InstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *InstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.InstanceId, ok = input.Parsed["instanceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceId", input)
	}

	return nil
}

// ValidateInstanceID checks that 'input' can be parsed as a Instance ID
func ValidateInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Instance ID
func (id InstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/instances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.InstanceId)
}

// Segments returns a slice of Resource ID Segments which comprise this Instance ID
func (id InstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceId", "instanceId"),
	}
}

// String returns a human-readable description of this Instance ID
func (id InstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Instance: %q", id.InstanceId),
	}
	return fmt.Sprintf("Instance (%s)", strings.Join(components, "\n"))
}
