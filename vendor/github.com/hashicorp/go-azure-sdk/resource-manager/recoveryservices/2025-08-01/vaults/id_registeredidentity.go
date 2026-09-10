package vaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RegisteredIdentityId{})
}

var _ resourceids.ResourceId = &RegisteredIdentityId{}

// RegisteredIdentityId is a struct representing the Resource ID for a Registered Identity
type RegisteredIdentityId struct {
	SubscriptionId         string
	ResourceGroupName      string
	VaultName              string
	RegisteredIdentityName string
}

// NewRegisteredIdentityID returns a new RegisteredIdentityId struct
func NewRegisteredIdentityID(subscriptionId string, resourceGroupName string, vaultName string, registeredIdentityName string) RegisteredIdentityId {
	return RegisteredIdentityId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		VaultName:              vaultName,
		RegisteredIdentityName: registeredIdentityName,
	}
}

// ParseRegisteredIdentityID parses 'input' into a RegisteredIdentityId
func ParseRegisteredIdentityID(input string) (*RegisteredIdentityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RegisteredIdentityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RegisteredIdentityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRegisteredIdentityIDInsensitively parses 'input' case-insensitively into a RegisteredIdentityId
// note: this method should only be used for API response data and not user input
func ParseRegisteredIdentityIDInsensitively(input string) (*RegisteredIdentityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RegisteredIdentityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RegisteredIdentityId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RegisteredIdentityId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VaultName, ok = input.Parsed["vaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vaultName", input)
	}

	if id.RegisteredIdentityName, ok = input.Parsed["registeredIdentityName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registeredIdentityName", input)
	}

	return nil
}

// ValidateRegisteredIdentityID checks that 'input' can be parsed as a Registered Identity ID
func ValidateRegisteredIdentityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRegisteredIdentityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Registered Identity ID
func (id RegisteredIdentityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.RecoveryServices/vaults/%s/registeredIdentities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.RegisteredIdentityName)
}

// Segments returns a slice of Resource ID Segments which comprise this Registered Identity ID
func (id RegisteredIdentityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices", "Microsoft.RecoveryServices"),
		resourceids.StaticSegment("staticVaults", "vaults", "vaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultName"),
		resourceids.StaticSegment("staticRegisteredIdentities", "registeredIdentities", "registeredIdentities"),
		resourceids.UserSpecifiedSegment("registeredIdentityName", "registeredIdentityName"),
	}
}

// String returns a human-readable description of this Registered Identity ID
func (id RegisteredIdentityId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Registered Identity Name: %q", id.RegisteredIdentityName),
	}
	return fmt.Sprintf("Registered Identity (%s)", strings.Join(components, "\n"))
}
