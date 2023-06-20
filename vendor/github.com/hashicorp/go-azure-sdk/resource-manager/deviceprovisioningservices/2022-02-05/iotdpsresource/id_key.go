package iotdpsresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = KeyId{}

// KeyId is a struct representing the Resource ID for a Key
type KeyId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ProvisioningServiceName string
	KeyName                 string
}

// NewKeyID returns a new KeyId struct
func NewKeyID(subscriptionId string, resourceGroupName string, provisioningServiceName string, keyName string) KeyId {
	return KeyId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ProvisioningServiceName: provisioningServiceName,
		KeyName:                 keyName,
	}
}

// ParseKeyID parses 'input' into a KeyId
func ParseKeyID(input string) (*KeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(KeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ProvisioningServiceName, ok = parsed.Parsed["provisioningServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "provisioningServiceName", *parsed)
	}

	if id.KeyName, ok = parsed.Parsed["keyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "keyName", *parsed)
	}

	return &id, nil
}

// ParseKeyIDInsensitively parses 'input' case-insensitively into a KeyId
// note: this method should only be used for API response data and not user input
func ParseKeyIDInsensitively(input string) (*KeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(KeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ProvisioningServiceName, ok = parsed.Parsed["provisioningServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "provisioningServiceName", *parsed)
	}

	if id.KeyName, ok = parsed.Parsed["keyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "keyName", *parsed)
	}

	return &id, nil
}

// ValidateKeyID checks that 'input' can be parsed as a Key ID
func ValidateKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Key ID
func (id KeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/provisioningServices/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProvisioningServiceName, id.KeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Key ID
func (id KeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevices", "Microsoft.Devices", "Microsoft.Devices"),
		resourceids.StaticSegment("staticProvisioningServices", "provisioningServices", "provisioningServices"),
		resourceids.UserSpecifiedSegment("provisioningServiceName", "provisioningServiceValue"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyValue"),
	}
}

// String returns a human-readable description of this Key ID
func (id KeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provisioning Service Name: %q", id.ProvisioningServiceName),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Key (%s)", strings.Join(components, "\n"))
}
