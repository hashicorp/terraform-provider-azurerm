package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DashboardsClient           *dashboard.DashboardClient
	TenantConfigurationsClient *tenantconfiguration.TenantConfigurationClient
}

func NewClient(o *common.ClientOptions) *Client {
	dashboardsClient := dashboard.NewDashboardClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&dashboardsClient.Client, o.ResourceManagerAuthorizer)

	tenantConfigurationsClient := tenantconfiguration.NewTenantConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&tenantConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DashboardsClient:           &dashboardsClient,
		TenantConfigurationsClient: &tenantConfigurationsClient,
	}
}
