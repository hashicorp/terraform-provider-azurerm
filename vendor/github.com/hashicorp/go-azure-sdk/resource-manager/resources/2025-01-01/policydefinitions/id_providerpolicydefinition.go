package policydefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProviderPolicyDefinitionId{})
}

var _ resourceids.ResourceId = &ProviderPolicyDefinitionId{}

// ProviderPolicyDefinitionId is a struct representing the Resource ID for a Provider Policy Definition
type ProviderPolicyDefinitionId struct {
	SubscriptionId       string
	PolicyDefinitionName string
}

// NewProviderPolicyDefinitionID returns a new ProviderPolicyDefinitionId struct
func NewProviderPolicyDefinitionID(subscriptionId string, policyDefinitionName string) ProviderPolicyDefinitionId {
	return ProviderPolicyDefinitionId{
		SubscriptionId:       subscriptionId,
		PolicyDefinitionName: policyDefinitionName,
	}
}

// ParseProviderPolicyDefinitionID parses 'input' into a ProviderPolicyDefinitionId
func ParseProviderPolicyDefinitionID(input string) (*ProviderPolicyDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderPolicyDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderPolicyDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProviderPolicyDefinitionIDInsensitively parses 'input' case-insensitively into a ProviderPolicyDefinitionId
// note: this method should only be used for API response data and not user input
func ParseProviderPolicyDefinitionIDInsensitively(input string) (*ProviderPolicyDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProviderPolicyDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProviderPolicyDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProviderPolicyDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.PolicyDefinitionName, ok = input.Parsed["policyDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policyDefinitionName", input)
	}

	return nil
}

// ValidateProviderPolicyDefinitionID checks that 'input' can be parsed as a Provider Policy Definition ID
func ValidateProviderPolicyDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderPolicyDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Policy Definition ID
func (id ProviderPolicyDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Authorization/policyDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PolicyDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Policy Definition ID
func (id ProviderPolicyDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAuthorization", "Microsoft.Authorization", "Microsoft.Authorization"),
		resourceids.StaticSegment("staticPolicyDefinitions", "policyDefinitions", "policyDefinitions"),
		resourceids.UserSpecifiedSegment("policyDefinitionName", "policyDefinitionName"),
	}
}

// String returns a human-readable description of this Provider Policy Definition ID
func (id ProviderPolicyDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Policy Definition Name: %q", id.PolicyDefinitionName),
	}
	return fmt.Sprintf("Provider Policy Definition (%s)", strings.Join(components, "\n"))
}
