package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AttestationProvidersId struct {
	SubscriptionId          string
	ResourceGroup           string
	AttestationProviderName string
}

func NewAttestationProvidersID(subscriptionId, resourceGroup, attestationProviderName string) AttestationProvidersId {
	return AttestationProvidersId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		AttestationProviderName: attestationProviderName,
	}
}

func (id AttestationProvidersId) String() string {
	segments := []string{
		fmt.Sprintf("Attestation Provider Name %q", id.AttestationProviderName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Attestation Providers", segmentsStr)
}

func (id AttestationProvidersId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Attestation/attestationProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AttestationProviderName)
}

// ParseAttestationProvidersID parses a AttestationProviders ID into an AttestationProvidersId struct
func ParseAttestationProvidersID(input string) (*AttestationProvidersId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AttestationProvidersId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AttestationProviderName, err = id.PopSegment("attestationProviders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseAttestationProvidersIDInsensitively parses an AttestationProviders ID into an AttestationProvidersId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseAttestationProvidersID method should be used instead for validation etc.
func ParseAttestationProvidersIDInsensitively(input string) (*AttestationProvidersId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AttestationProvidersId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'attestationProviders' segment
	attestationProvidersKey := "attestationProviders"
	for key := range id.Path {
		if strings.EqualFold(key, attestationProvidersKey) {
			attestationProvidersKey = key
			break
		}
	}
	if resourceId.AttestationProviderName, err = id.PopSegment(attestationProvidersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
