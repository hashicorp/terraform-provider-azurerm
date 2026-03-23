package tokens

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TokenId{})
}

var _ resourceids.ResourceId = &TokenId{}

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
	parser := resourceids.NewParserFromResourceIdType(&TokenId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TokenId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTokenIDInsensitively parses 'input' case-insensitively into a TokenId
// note: this method should only be used for API response data and not user input
func ParseTokenIDInsensitively(input string) (*TokenId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TokenId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TokenId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TokenId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.RegistryName, ok = input.Parsed["registryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "registryName", input)
	}

	if id.TokenName, ok = input.Parsed["tokenName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tokenName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("registryName", "registryName"),
		resourceids.StaticSegment("staticTokens", "tokens", "tokens"),
		resourceids.UserSpecifiedSegment("tokenName", "tokenName"),
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
