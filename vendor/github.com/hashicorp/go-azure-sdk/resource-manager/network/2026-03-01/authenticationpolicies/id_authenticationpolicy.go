package authenticationpolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AuthenticationPolicyId{})
}

var _ resourceids.ResourceId = &AuthenticationPolicyId{}

// AuthenticationPolicyId is a struct representing the Resource ID for a Authentication Policy
type AuthenticationPolicyId struct {
	SubscriptionId           string
	ResourceGroupName        string
	AuthenticationPolicyName string
}

// NewAuthenticationPolicyID returns a new AuthenticationPolicyId struct
func NewAuthenticationPolicyID(subscriptionId string, resourceGroupName string, authenticationPolicyName string) AuthenticationPolicyId {
	return AuthenticationPolicyId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		AuthenticationPolicyName: authenticationPolicyName,
	}
}

// ParseAuthenticationPolicyID parses 'input' into a AuthenticationPolicyId
func ParseAuthenticationPolicyID(input string) (*AuthenticationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthenticationPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthenticationPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAuthenticationPolicyIDInsensitively parses 'input' case-insensitively into a AuthenticationPolicyId
// note: this method should only be used for API response data and not user input
func ParseAuthenticationPolicyIDInsensitively(input string) (*AuthenticationPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AuthenticationPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AuthenticationPolicyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AuthenticationPolicyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AuthenticationPolicyName, ok = input.Parsed["authenticationPolicyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "authenticationPolicyName", input)
	}

	return nil
}

// ValidateAuthenticationPolicyID checks that 'input' can be parsed as a Authentication Policy ID
func ValidateAuthenticationPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAuthenticationPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Authentication Policy ID
func (id AuthenticationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/authenticationPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AuthenticationPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Authentication Policy ID
func (id AuthenticationPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticAuthenticationPolicies", "authenticationPolicies", "authenticationPolicies"),
		resourceids.UserSpecifiedSegment("authenticationPolicyName", "authenticationPolicyName"),
	}
}

// String returns a human-readable description of this Authentication Policy ID
func (id AuthenticationPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Authentication Policy Name: %q", id.AuthenticationPolicyName),
	}
	return fmt.Sprintf("Authentication Policy (%s)", strings.Join(components, "\n"))
}
