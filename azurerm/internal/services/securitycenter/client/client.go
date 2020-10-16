package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ContactsClient                 *security.ContactsClient
	PricingClient                  *security.PricingsClient
	WorkspaceClient                *security.WorkspaceSettingsClient
	AdvancedThreatProtectionClient *security.AdvancedThreatProtectionClient
	SettingClient                  *security.SettingsClient
	// Note. This is patched version of the client
	// - With fixes for https://github.com/Azure/azure-sdk-for-go/issues/12634
	AutomationsClient *azuresdkhacks.AutomationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ascLocation := "Global"

	ContactsClient := security.NewContactsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ContactsClient.Client, o.ResourceManagerAuthorizer)

	PricingClient := security.NewPricingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&PricingClient.Client, o.ResourceManagerAuthorizer)

	WorkspaceClient := security.NewWorkspaceSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&WorkspaceClient.Client, o.ResourceManagerAuthorizer)

	AdvancedThreatProtectionClient := security.NewAdvancedThreatProtectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AdvancedThreatProtectionClient.Client, o.ResourceManagerAuthorizer)

	SettingClient := security.NewSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&SettingClient.Client, o.ResourceManagerAuthorizer)
	AutomationsClient := azuresdkhacks.NewAutomationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AutomationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ContactsClient:                 &ContactsClient,
		PricingClient:                  &PricingClient,
		WorkspaceClient:                &WorkspaceClient,
		AdvancedThreatProtectionClient: &AdvancedThreatProtectionClient,
		SettingClient:                  &SettingClient,
		AutomationsClient:              &AutomationsClient,
	}
}
