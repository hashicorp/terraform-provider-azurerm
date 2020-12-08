package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ProviderId struct {
	SubscriptionId          string
	ResourceGroup           string
	AttestationProviderName string
}

func NewProviderID(subscriptionId, resourceGroup, attestationProviderName string) ProviderId {
	return ProviderId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		AttestationProviderName: attestationProviderName,
	}
}

func (id ProviderId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Attestation Provider Name %q", id.AttestationProviderName),
	}
	return strings.Join(segments, " / ")
}

func (id ProviderId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Attestation/attestationProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AttestationProviderName)
}

// ProviderID parses a Provider ID into an ProviderId struct
func ProviderID(input string) (*ProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ProviderId{
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
