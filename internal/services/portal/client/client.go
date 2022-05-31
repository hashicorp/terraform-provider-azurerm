package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/sdk/2019-01-01-preview/dashboard"
)

type Client struct {
	DashboardsClient           *dashboard.DashboardClient
	TenantConfigurationsClient *portal.TenantConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	dashboardsClient := dashboard.NewDashboardClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dashboardsClient.Client, o.ResourceManagerAuthorizer)

	tenantConfigurationsClient := portal.NewTenantConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tenantConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DashboardsClient:           &dashboardsClient,
		TenantConfigurationsClient: &tenantConfigurationsClient,
	}
}
