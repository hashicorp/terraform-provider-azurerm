package attestationproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AttestationProvidersId{}

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
	parser := resourceids.NewParserFromResourceIdType(AttestationProvidersId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AttestationProvidersId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AttestationProviderName, ok = parsed.Parsed["attestationProviderName"]; !ok {
		return nil, fmt.Errorf("the segment 'attestationProviderName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAttestationProvidersIDInsensitively parses 'input' case-insensitively into a AttestationProvidersId
// note: this method should only be used for API response data and not user input
func ParseAttestationProvidersIDInsensitively(input string) (*AttestationProvidersId, error) {
	parser := resourceids.NewParserFromResourceIdType(AttestationProvidersId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AttestationProvidersId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AttestationProviderName, ok = parsed.Parsed["attestationProviderName"]; !ok {
		return nil, fmt.Errorf("the segment 'attestationProviderName' was not found in the resource id %q", input)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("attestationProviderName", "attestationProviderValue"),
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
