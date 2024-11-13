package providerinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderInstanceId{})
}

var _ resourceids.ResourceId = &ProviderInstanceId{}

// ProviderInstanceId is a struct representing the Resource ID for a Provider Instance
type ProviderInstanceId struct {
	SubscriptionId       string
	ResourceGroupName    string
	MonitorName          string
	ProviderInstanceName string
}

// NewProviderInstanceID returns a new ProviderInstanceId struct
func NewProviderInstanceID(subscriptionId string, resourceGroupName string, monitorName string, providerInstanceName string) ProviderInstanceId {
	return ProviderInstanceId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		MonitorName:          monitorName,
		ProviderInstanceName: providerInstanceName,
	}
}

// ParseProviderInstanceID parses 'input' into a ProviderInstanceId
func ParseProviderInstanceID(input string) (*ProviderInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderInstanceIDInsensitively parses 'input' case-insensitively into a ProviderInstanceId
// note: this method should only be used for API response data and not user input
func ParseProviderInstanceIDInsensitively(input string) (*ProviderInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MonitorName, ok = input.Parsed["monitorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "monitorName", input)
	}

	if id.ProviderInstanceName, ok = input.Parsed["providerInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "providerInstanceName", input)
	}

	return nil
}

// ValidateProviderInstanceID checks that 'input' can be parsed as a Provider Instance ID
func ValidateProviderInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Instance ID
func (id ProviderInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Workloads/monitors/%s/providerInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName, id.ProviderInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Instance ID
func (id ProviderInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWorkloads", "Microsoft.Workloads", "Microsoft.Workloads"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorName"),
		resourceids.StaticSegment("staticProviderInstances", "providerInstances", "providerInstances"),
		resourceids.UserSpecifiedSegment("providerInstanceName", "providerInstanceName"),
	}
}

// String returns a human-readable description of this Provider Instance ID
func (id ProviderInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
		fmt.Sprintf("Provider Instance Name: %q", id.ProviderInstanceName),
	}
	return fmt.Sprintf("Provider Instance (%s)", strings.Join(components, "\n"))
}
