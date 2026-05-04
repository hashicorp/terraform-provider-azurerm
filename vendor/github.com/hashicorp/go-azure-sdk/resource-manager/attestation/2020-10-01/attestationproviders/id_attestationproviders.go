package attestationproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AttestationProvidersId{})
}

var _ resourceids.ResourceId = &AttestationProvidersId{}

// AttestationProvidersId is a struct representing the Resource ID for a Attestation Providers
type AttestationProvidersId struct {
	SubscriptionId          string
	ResourceGroupName       string
	AttestationProviderName string
}

// NewAttestationProvidersID returns a new AttestationProvidersId struct
func NewAttestationProvidersID(subscriptionId string, resourceGroupName string, attestationProviderName string) AttestationProvidersId {
	return AttestationProvidersId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		AttestationProviderName: attestationProviderName,
	}
}

// ParseAttestationProvidersID parses 'input' into a AttestationProvidersId
func ParseAttestationProvidersID(input string) (*AttestationProvidersId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AttestationProvidersId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AttestationProvidersId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAttestationProvidersIDInsensitively parses 'input' case-insensitively into a AttestationProvidersId
// note: this method should only be used for API response data and not user input
func ParseAttestationProvidersIDInsensitively(input string) (*AttestationProvidersId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AttestationProvidersId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AttestationProvidersId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AttestationProvidersId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AttestationProviderName, ok = input.Parsed["attestationProviderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "attestationProviderName", input)
	}

	return nil
}

// ValidateAttestationProvidersID checks that 'input' can be parsed as a Attestation Providers ID
func ValidateAttestationProvidersID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAttestationProvidersID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Attestation Providers ID
func (id AttestationProvidersId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Attestation/attestationProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AttestationProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Attestation Providers ID
func (id AttestationProvidersId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAttestation", "Microsoft.Attestation", "Microsoft.Attestation"),
		resourceids.StaticSegment("staticAttestationProviders", "attestationProviders", "attestationProviders"),
		resourceids.UserSpecifiedSegment("attestationProviderName", "attestationProviderName"),
	}
}

// String returns a human-readable description of this Attestation Providers ID
func (id AttestationProvidersId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Attestation Provider Name: %q", id.AttestationProviderName),
	}
	return fmt.Sprintf("Attestation Providers (%s)", strings.Join(components, "\n"))
}
