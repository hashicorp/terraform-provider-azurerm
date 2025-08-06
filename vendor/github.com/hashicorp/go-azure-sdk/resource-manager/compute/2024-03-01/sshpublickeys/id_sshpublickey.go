package sshpublickeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SshPublicKeyId{})
}

var _ resourceids.ResourceId = &SshPublicKeyId{}

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
	parser := resourceids.NewParserFromResourceIdType(&SshPublicKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SshPublicKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSshPublicKeyIDInsensitively parses 'input' case-insensitively into a SshPublicKeyId
// note: this method should only be used for API response data and not user input
func ParseSshPublicKeyIDInsensitively(input string) (*SshPublicKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SshPublicKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SshPublicKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SshPublicKeyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SshPublicKeyName, ok = input.Parsed["sshPublicKeyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sshPublicKeyName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("sshPublicKeyName", "sshPublicKeyName"),
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
