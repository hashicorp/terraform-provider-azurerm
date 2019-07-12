package securitycenter

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ContactsClient                 security.ContactsClient
	PricingClient                  security.PricingsClient
	WorkspaceClient                security.WorkspaceSettingsClient
	AdvancedThreatProtectionClient security.AdvancedThreatProtectionClient
}

func BuildClient(o *common.ClientOptions) *Client {
	ascLocation := "Global"
	c := Client{}

	c.ContactsClient = security.NewContactsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&c.ContactsClient.Client, o.ResourceManagerAuthorizer)

	c.PricingClient = security.NewPricingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&c.PricingClient.Client, o.ResourceManagerAuthorizer)

	c.WorkspaceClient = security.NewWorkspaceSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&c.WorkspaceClient.Client, o.ResourceManagerAuthorizer)

	c.AdvancedThreatProtectionClient = security.NewAdvancedThreatProtectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&c.AdvancedThreatProtectionClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
