package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApiManagementIdentityProviderId struct {
	ResourceGroup string
	ServiceName   string
	ProviderName  string
}

func (a *ApiManagementIdentityProviderId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/identityProviders/%s",
		subscriptionId, a.ResourceGroup, a.ServiceName, a.ProviderName)
}

func (a *ApiManagementIdentityProviderId) ApiManagementID(subscriptionId string) string {
	id := ApiManagementId{
		ResourceGroup: a.ResourceGroup,
		ServiceName:   a.ServiceName,
	}
	return id.ID(subscriptionId)
}

func ApiManagementIdentityProviderID(input string) (*ApiManagementIdentityProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	apiManagementIdentityProvider := ApiManagementIdentityProviderId{
		ResourceGroup: id.ResourceGroup,
	}

	if apiManagementIdentityProvider.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}

	if apiManagementIdentityProvider.ProviderName, err = id.PopSegment("identityProviders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &apiManagementIdentityProvider, nil
}
