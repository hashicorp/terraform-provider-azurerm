package deviceupdates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = InstanceId{}

// InstanceId is a struct representing the Resource ID for a Instance
type InstanceId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	InstanceName      string
}

// NewInstanceID returns a new InstanceId struct
func NewInstanceID(subscriptionId string, resourceGroupName string, accountName string, instanceName string) InstanceId {
	return InstanceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		InstanceName:      instanceName,
	}
}

// ParseInstanceID parses 'input' into a InstanceId
func ParseInstanceID(input string) (*InstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(InstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.InstanceName, ok = parsed.Parsed["instanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "instanceName", *parsed)
	}

	return &id, nil
}

// ParseInstanceIDInsensitively parses 'input' case-insensitively into a InstanceId
// note: this method should only be used for API response data and not user input
func ParseInstanceIDInsensitively(input string) (*InstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(InstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "accountName", *parsed)
	}

	if id.InstanceName, ok = parsed.Parsed["instanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "instanceName", *parsed)
	}

	return &id, nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DeviceUpdate/accounts/%s/instances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.InstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Instance ID
func (id InstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDeviceUpdate", "Microsoft.DeviceUpdate", "Microsoft.DeviceUpdate"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceValue"),
	}
}

// String returns a human-readable description of this Instance ID
func (id InstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
	}
	return fmt.Sprintf("Instance (%s)", strings.Join(components, "\n"))
}
