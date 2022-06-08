package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DashboardsClient           *dashboard.DashboardClient
	LegacyDashboardsClient     *portal.DashboardsClient
	TenantConfigurationsClient *tenantconfiguration.TenantConfigurationClient
}

func NewClient(o *common.ClientOptions) *Client {
	dashboardsClient := dashboard.NewDashboardClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dashboardsClient.Client, o.ResourceManagerAuthorizer)

	legacyDashboardsClient := portal.NewDashboardsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyDashboardsClient.Client, o.ResourceManagerAuthorizer)

	tenantConfigurationsClient := tenantconfiguration.NewTenantConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&tenantConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DashboardsClient:           &dashboardsClient,
		LegacyDashboardsClient:     &legacyDashboardsClient,
		TenantConfigurationsClient: &tenantConfigurationsClient,
	}
}
