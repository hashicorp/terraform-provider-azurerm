package client

import (
	securityv1 "github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	securityv3 "github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ContactsClient                 *securityv1.ContactsClient
	PricingClient                  *securityv3.PricingsClient
	WorkspaceClient                *securityv1.WorkspaceSettingsClient
	AdvancedThreatProtectionClient *securityv1.AdvancedThreatProtectionClient
}

func NewClient(o *common.ClientOptions) *Client {
	ascLocation := "Global"

	ContactsClient := securityv1.NewContactsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&ContactsClient.Client, o.ResourceManagerAuthorizer)

	PricingClient := securityv3.NewPricingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&PricingClient.Client, o.ResourceManagerAuthorizer)

	WorkspaceClient := securityv1.NewWorkspaceSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&WorkspaceClient.Client, o.ResourceManagerAuthorizer)

	AdvancedThreatProtectionClient := securityv1.NewAdvancedThreatProtectionClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId, ascLocation)
	o.ConfigureClient(&AdvancedThreatProtectionClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ContactsClient:                 &ContactsClient,
		PricingClient:                  &PricingClient,
		WorkspaceClient:                &WorkspaceClient,
		AdvancedThreatProtectionClient: &AdvancedThreatProtectionClient,
	}
}
