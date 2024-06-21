package managedclusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GuardrailsVersionId{})
}

var _ resourceids.ResourceId = &GuardrailsVersionId{}

// GuardrailsVersionId is a struct representing the Resource ID for a Guardrails Version
type GuardrailsVersionId struct {
	SubscriptionId        string
	LocationName          string
	GuardrailsVersionName string
}

// NewGuardrailsVersionID returns a new GuardrailsVersionId struct
func NewGuardrailsVersionID(subscriptionId string, locationName string, guardrailsVersionName string) GuardrailsVersionId {
	return GuardrailsVersionId{
		SubscriptionId:        subscriptionId,
		LocationName:          locationName,
		GuardrailsVersionName: guardrailsVersionName,
	}
}

// ParseGuardrailsVersionID parses 'input' into a GuardrailsVersionId
func ParseGuardrailsVersionID(input string) (*GuardrailsVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GuardrailsVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GuardrailsVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGuardrailsVersionIDInsensitively parses 'input' case-insensitively into a GuardrailsVersionId
// note: this method should only be used for API response data and not user input
func ParseGuardrailsVersionIDInsensitively(input string) (*GuardrailsVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GuardrailsVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GuardrailsVersionId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GuardrailsVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.GuardrailsVersionName, ok = input.Parsed["guardrailsVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "guardrailsVersionName", input)
	}

	return nil
}

// ValidateGuardrailsVersionID checks that 'input' can be parsed as a Guardrails Version ID
func ValidateGuardrailsVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGuardrailsVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Guardrails Version ID
func (id GuardrailsVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.ContainerService/locations/%s/guardrailsVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.GuardrailsVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Guardrails Version ID
func (id GuardrailsVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticGuardrailsVersions", "guardrailsVersions", "guardrailsVersions"),
		resourceids.UserSpecifiedSegment("guardrailsVersionName", "guardrailsVersionValue"),
	}
}

// String returns a human-readable description of this Guardrails Version ID
func (id GuardrailsVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Guardrails Version Name: %q", id.GuardrailsVersionName),
	}
	return fmt.Sprintf("Guardrails Version (%s)", strings.Join(components, "\n"))
}
