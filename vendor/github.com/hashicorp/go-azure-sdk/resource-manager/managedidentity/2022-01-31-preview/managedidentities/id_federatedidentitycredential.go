package managedidentities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FederatedIdentityCredentialId{}

// FederatedIdentityCredentialId is a struct representing the Resource ID for a Federated Identity Credential
type FederatedIdentityCredentialId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	UserAssignedIdentityName        string
	FederatedIdentityCredentialName string
}

// NewFederatedIdentityCredentialID returns a new FederatedIdentityCredentialId struct
func NewFederatedIdentityCredentialID(subscriptionId string, resourceGroupName string, userAssignedIdentityName string, federatedIdentityCredentialName string) FederatedIdentityCredentialId {
	return FederatedIdentityCredentialId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		UserAssignedIdentityName:        userAssignedIdentityName,
		FederatedIdentityCredentialName: federatedIdentityCredentialName,
	}
}

// ParseFederatedIdentityCredentialID parses 'input' into a FederatedIdentityCredentialId
func ParseFederatedIdentityCredentialID(input string) (*FederatedIdentityCredentialId, error) {
	parser := resourceids.NewParserFromResourceIdType(FederatedIdentityCredentialId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FederatedIdentityCredentialId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.UserAssignedIdentityName, ok = parsed.Parsed["userAssignedIdentityName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userAssignedIdentityName", *parsed)
	}

	if id.FederatedIdentityCredentialName, ok = parsed.Parsed["federatedIdentityCredentialName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "federatedIdentityCredentialName", *parsed)
	}

	return &id, nil
}

// ParseFederatedIdentityCredentialIDInsensitively parses 'input' case-insensitively into a FederatedIdentityCredentialId
// note: this method should only be used for API response data and not user input
func ParseFederatedIdentityCredentialIDInsensitively(input string) (*FederatedIdentityCredentialId, error) {
	parser := resourceids.NewParserFromResourceIdType(FederatedIdentityCredentialId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FederatedIdentityCredentialId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.UserAssignedIdentityName, ok = parsed.Parsed["userAssignedIdentityName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "userAssignedIdentityName", *parsed)
	}

	if id.FederatedIdentityCredentialName, ok = parsed.Parsed["federatedIdentityCredentialName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "federatedIdentityCredentialName", *parsed)
	}

	return &id, nil
}

// ValidateFederatedIdentityCredentialID checks that 'input' can be parsed as a Federated Identity Credential ID
func ValidateFederatedIdentityCredentialID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFederatedIdentityCredentialID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Federated Identity Credential ID
func (id FederatedIdentityCredentialId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s/federatedIdentityCredentials/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.UserAssignedIdentityName, id.FederatedIdentityCredentialName)
}

// Segments returns a slice of Resource ID Segments which comprise this Federated Identity Credential ID
func (id FederatedIdentityCredentialId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftManagedIdentity", "Microsoft.ManagedIdentity", "Microsoft.ManagedIdentity"),
		resourceids.StaticSegment("staticUserAssignedIdentities", "userAssignedIdentities", "userAssignedIdentities"),
		resourceids.UserSpecifiedSegment("userAssignedIdentityName", "userAssignedIdentityValue"),
		resourceids.StaticSegment("staticFederatedIdentityCredentials", "federatedIdentityCredentials", "federatedIdentityCredentials"),
		resourceids.UserSpecifiedSegment("federatedIdentityCredentialName", "federatedIdentityCredentialValue"),
	}
}

// String returns a human-readable description of this Federated Identity Credential ID
func (id FederatedIdentityCredentialId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("User Assigned Identity Name: %q", id.UserAssignedIdentityName),
		fmt.Sprintf("Federated Identity Credential Name: %q", id.FederatedIdentityCredentialName),
	}
	return fmt.Sprintf("Federated Identity Credential (%s)", strings.Join(components, "\n"))
}
