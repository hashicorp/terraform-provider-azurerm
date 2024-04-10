package credentials

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &CredentialId{}

// CredentialId is a struct representing the Resource ID for a Credential
type CredentialId struct {
	SubscriptionId    string
	ResourceGroupName string
	FactoryName       string
	CredentialName    string
}

// NewCredentialID returns a new CredentialId struct
func NewCredentialID(subscriptionId string, resourceGroupName string, factoryName string, credentialName string) CredentialId {
	return CredentialId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FactoryName:       factoryName,
		CredentialName:    credentialName,
	}
}

// ParseCredentialID parses 'input' into a CredentialId
func ParseCredentialID(input string) (*CredentialId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CredentialId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CredentialId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCredentialIDInsensitively parses 'input' case-insensitively into a CredentialId
// note: this method should only be used for API response data and not user input
func ParseCredentialIDInsensitively(input string) (*CredentialId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CredentialId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CredentialId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CredentialId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FactoryName, ok = input.Parsed["factoryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "factoryName", input)
	}

	if id.CredentialName, ok = input.Parsed["credentialName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "credentialName", input)
	}

	return nil
}

// ValidateCredentialID checks that 'input' can be parsed as a Credential ID
func ValidateCredentialID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCredentialID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Credential ID
func (id CredentialId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s/credentials/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.CredentialName)
}

// Segments returns a slice of Resource ID Segments which comprise this Credential ID
func (id CredentialId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataFactory", "Microsoft.DataFactory", "Microsoft.DataFactory"),
		resourceids.StaticSegment("staticFactories", "factories", "factories"),
		resourceids.UserSpecifiedSegment("factoryName", "factoryValue"),
		resourceids.StaticSegment("staticCredentials", "credentials", "credentials"),
		resourceids.UserSpecifiedSegment("credentialName", "credentialValue"),
	}
}

// String returns a human-readable description of this Credential ID
func (id CredentialId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Factory Name: %q", id.FactoryName),
		fmt.Sprintf("Credential Name: %q", id.CredentialName),
	}
	return fmt.Sprintf("Credential (%s)", strings.Join(components, "\n"))
}
