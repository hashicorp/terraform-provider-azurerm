package policysetdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderPolicySetDefinitionId{})
}

var _ resourceids.ResourceId = &ProviderPolicySetDefinitionId{}

// ProviderPolicySetDefinitionId is a struct representing the Resource ID for a Provider Policy Set Definition
type ProviderPolicySetDefinitionId struct {
	SubscriptionId          string
	PolicySetDefinitionName string
}

// NewProviderPolicySetDefinitionID returns a new ProviderPolicySetDefinitionId struct
func NewProviderPolicySetDefinitionID(subscriptionId string, policySetDefinitionName string) ProviderPolicySetDefinitionId {
	return ProviderPolicySetDefinitionId{
		SubscriptionId:          subscriptionId,
		PolicySetDefinitionName: policySetDefinitionName,
	}
}

// ParseProviderPolicySetDefinitionID parses 'input' into a ProviderPolicySetDefinitionId
func ParseProviderPolicySetDefinitionID(input string) (*ProviderPolicySetDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderPolicySetDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderPolicySetDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderPolicySetDefinitionIDInsensitively parses 'input' case-insensitively into a ProviderPolicySetDefinitionId
// note: this method should only be used for API response data and not user input
func ParseProviderPolicySetDefinitionIDInsensitively(input string) (*ProviderPolicySetDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderPolicySetDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderPolicySetDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderPolicySetDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.PolicySetDefinitionName, ok = input.Parsed["policySetDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policySetDefinitionName", input)
	}

	return nil
}

// ValidateProviderPolicySetDefinitionID checks that 'input' can be parsed as a Provider Policy Set Definition ID
func ValidateProviderPolicySetDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderPolicySetDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Policy Set Definition ID
func (id ProviderPolicySetDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Authorization/policySetDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PolicySetDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Policy Set Definition ID
func (id ProviderPolicySetDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticPolicySetDefinitions", "policySetDefinitions", "policySetDefinitions"),
		resourceids.UserSpecifiedSegment("policySetDefinitionName", "policySetDefinitionName"),
	}
}

// String returns a human-readable description of this Provider Policy Set Definition ID
func (id ProviderPolicySetDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Policy Set Definition Name: %q", id.PolicySetDefinitionName),
	}
	return fmt.Sprintf("Provider Policy Set Definition (%s)", strings.Join(components, "\n"))
}
