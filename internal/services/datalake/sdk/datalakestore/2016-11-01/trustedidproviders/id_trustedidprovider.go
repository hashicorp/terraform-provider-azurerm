package trustedidproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type TrustedIdProviderId struct {
	SubscriptionId string
	ResourceGroup  string
	AccountName    string
	Name           string
}

func NewTrustedIdProviderID(subscriptionId, resourceGroup, accountName, name string) TrustedIdProviderId {
	return TrustedIdProviderId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		AccountName:    accountName,
		Name:           name,
	}
}

func (id TrustedIdProviderId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Account Name %q", id.AccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Trusted Id Provider", segmentsStr)
}

func (id TrustedIdProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeStore/accounts/%s/trustedIdProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AccountName, id.Name)
}

// ParseTrustedIdProviderID parses a TrustedIdProvider ID into an TrustedIdProviderId struct
func ParseTrustedIdProviderID(input string) (*TrustedIdProviderId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TrustedIdProviderId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("trustedIdProviders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseTrustedIdProviderIDInsensitively parses an TrustedIdProvider ID into an TrustedIdProviderId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseTrustedIdProviderID method should be used instead for validation etc.
func ParseTrustedIdProviderIDInsensitively(input string) (*TrustedIdProviderId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TrustedIdProviderId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'accounts' segment
	accountsKey := "accounts"
	for key := range id.Path {
		if strings.EqualFold(key, accountsKey) {
			accountsKey = key
			break
		}
	}
	if resourceId.AccountName, err = id.PopSegment(accountsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'trustedIdProviders' segment
	trustedIdProvidersKey := "trustedIdProviders"
	for key := range id.Path {
		if strings.EqualFold(key, trustedIdProvidersKey) {
			trustedIdProvidersKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(trustedIdProvidersKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
