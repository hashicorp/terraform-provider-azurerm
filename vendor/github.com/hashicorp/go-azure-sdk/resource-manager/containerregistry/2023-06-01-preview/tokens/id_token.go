package tokens

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TokenId{}

// TokenId is a struct representing the Resource ID for a Token
type TokenId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	TokenName         string
}

// NewTokenID returns a new TokenId struct
func NewTokenID(subscriptionId string, resourceGroupName string, registryName string, tokenName string) TokenId {
	return TokenId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		TokenName:         tokenName,
	}
}

// ParseTokenID parses 'input' into a TokenId
func ParseTokenID(input string) (*TokenId, error) {
	parser := resourceids.NewParserFromResourceIdType(TokenId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TokenId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.TokenName, ok = parsed.Parsed["tokenName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "tokenName", *parsed)
	}

	return &id, nil
}

// ParseTokenIDInsensitively parses 'input' case-insensitively into a TokenId
// note: this method should only be used for API response data and not user input
func ParseTokenIDInsensitively(input string) (*TokenId, error) {
	parser := resourceids.NewParserFromResourceIdType(TokenId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TokenId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.TokenName, ok = parsed.Parsed["tokenName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "tokenName", *parsed)
	}

	return &id, nil
}

// ValidateTokenID checks that 'input' can be parsed as a Token ID
func ValidateTokenID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTokenID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Token ID
func (id TokenId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/tokens/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.TokenName)
}

// Segments returns a slice of Resource ID Segments which comprise this Token ID
func (id TokenId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerRegistry", "Microsoft.ContainerRegistry", "Microsoft.ContainerRegistry"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticTokens", "tokens", "tokens"),
		resourceids.UserSpecifiedSegment("tokenName", "tokenValue"),
	}
}

// String returns a human-readable description of this Token ID
func (id TokenId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Token Name: %q", id.TokenName),
	}
	return fmt.Sprintf("Token (%s)", strings.Join(components, "\n"))
}
