package credentialsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CredentialSetId{}

// CredentialSetId is a struct representing the Resource ID for a Credential Set
type CredentialSetId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	CredentialSetName string
}

// NewCredentialSetID returns a new CredentialSetId struct
func NewCredentialSetID(subscriptionId string, resourceGroupName string, registryName string, credentialSetName string) CredentialSetId {
	return CredentialSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		CredentialSetName: credentialSetName,
	}
}

// ParseCredentialSetID parses 'input' into a CredentialSetId
func ParseCredentialSetID(input string) (*CredentialSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(CredentialSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CredentialSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.CredentialSetName, ok = parsed.Parsed["credentialSetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "credentialSetName", *parsed)
	}

	return &id, nil
}

// ParseCredentialSetIDInsensitively parses 'input' case-insensitively into a CredentialSetId
// note: this method should only be used for API response data and not user input
func ParseCredentialSetIDInsensitively(input string) (*CredentialSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(CredentialSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CredentialSetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.CredentialSetName, ok = parsed.Parsed["credentialSetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "credentialSetName", *parsed)
	}

	return &id, nil
}

// ValidateCredentialSetID checks that 'input' can be parsed as a Credential Set ID
func ValidateCredentialSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCredentialSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Credential Set ID
func (id CredentialSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/credentialSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.CredentialSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Credential Set ID
func (id CredentialSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticCredentialSets", "credentialSets", "credentialSets"),
		resourceids.UserSpecifiedSegment("credentialSetName", "credentialSetValue"),
	}
}

// String returns a human-readable description of this Credential Set ID
func (id CredentialSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Credential Set Name: %q", id.CredentialSetName),
	}
	return fmt.Sprintf("Credential Set (%s)", strings.Join(components, "\n"))
}
