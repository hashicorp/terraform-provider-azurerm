package sshpublickeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SshPublicKeyId{}

// SshPublicKeyId is a struct representing the Resource ID for a Ssh Public Key
type SshPublicKeyId struct {
	SubscriptionId    string
	ResourceGroupName string
	SshPublicKeyName  string
}

// NewSshPublicKeyID returns a new SshPublicKeyId struct
func NewSshPublicKeyID(subscriptionId string, resourceGroupName string, sshPublicKeyName string) SshPublicKeyId {
	return SshPublicKeyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SshPublicKeyName:  sshPublicKeyName,
	}
}

// ParseSshPublicKeyID parses 'input' into a SshPublicKeyId
func ParseSshPublicKeyID(input string) (*SshPublicKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(SshPublicKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SshPublicKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SshPublicKeyName, ok = parsed.Parsed["sshPublicKeyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sshPublicKeyName", *parsed)
	}

	return &id, nil
}

// ParseSshPublicKeyIDInsensitively parses 'input' case-insensitively into a SshPublicKeyId
// note: this method should only be used for API response data and not user input
func ParseSshPublicKeyIDInsensitively(input string) (*SshPublicKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(SshPublicKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SshPublicKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SshPublicKeyName, ok = parsed.Parsed["sshPublicKeyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sshPublicKeyName", *parsed)
	}

	return &id, nil
}

// ValidateSshPublicKeyID checks that 'input' can be parsed as a Ssh Public Key ID
func ValidateSshPublicKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSshPublicKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ssh Public Key ID
func (id SshPublicKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/sshPublicKeys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SshPublicKeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ssh Public Key ID
func (id SshPublicKeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticSshPublicKeys", "sshPublicKeys", "sshPublicKeys"),
		resourceids.UserSpecifiedSegment("sshPublicKeyName", "sshPublicKeyValue"),
	}
}

// String returns a human-readable description of this Ssh Public Key ID
func (id SshPublicKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ssh Public Key Name: %q", id.SshPublicKeyName),
	}
	return fmt.Sprintf("Ssh Public Key (%s)", strings.Join(components, "\n"))
}
