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
	recaser.RegisterResourceId(&DeletedVaultId{})
}

var _ resourceids.ResourceId = &DeletedVaultId{}

// DeletedVaultId is a struct representing the Resource ID for a Deleted Vault
type DeletedVaultId struct {
	SubscriptionId   string
	LocationName     string
	DeletedVaultName string
}

// NewDeletedVaultID returns a new DeletedVaultId struct
func NewDeletedVaultID(subscriptionId string, locationName string, deletedVaultName string) DeletedVaultId {
	return DeletedVaultId{
		SubscriptionId:   subscriptionId,
		LocationName:     locationName,
		DeletedVaultName: deletedVaultName,
	}
}

// ParseDeletedVaultID parses 'input' into a DeletedVaultId
func ParseDeletedVaultID(input string) (*DeletedVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedVaultId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedVaultIDInsensitively parses 'input' case-insensitively into a DeletedVaultId
// note: this method should only be used for API response data and not user input
func ParseDeletedVaultIDInsensitively(input string) (*DeletedVaultId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedVaultId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedVaultId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedVaultId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DeletedVaultName, ok = input.Parsed["deletedVaultName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedVaultName", input)
	}

	return nil
}

// ValidateDeletedVaultID checks that 'input' can be parsed as a Deleted Vault ID
func ValidateDeletedVaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deleted Vault ID
func (id DeletedVaultId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.KeyVault/locations/%s/deletedVaults/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DeletedVaultName)
}

// Segments returns a slice of Resource ID Segments which comprise this Deleted Vault ID
func (id DeletedVaultId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftKeyVault", "Microsoft.KeyVault", "Microsoft.KeyVault"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDeletedVaults", "deletedVaults", "deletedVaults"),
		resourceids.UserSpecifiedSegment("deletedVaultName", "deletedVaultName"),
	}
}

// String returns a human-readable description of this Deleted Vault ID
func (id DeletedVaultId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Deleted Vault Name: %q", id.DeletedVaultName),
	}
	return fmt.Sprintf("Deleted Vault (%s)", strings.Join(components, "\n"))
}
