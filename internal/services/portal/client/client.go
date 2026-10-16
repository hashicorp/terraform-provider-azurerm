// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/dashboards"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2026-04-01/tenantconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DashboardsClient           *dashboards.DashboardsClient
	TenantConfigurationsClient *tenantconfigurations.TenantConfigurationsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dashboardsClient, err := dashboards.NewDashboardsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dashboard client: %+v", err)
	}
	o.Configure(dashboardsClient.Client, o.Authorizers.ResourceManager)

	tenantConfigurationsClient, err := tenantconfigurations.NewTenantConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TenantConfiguration client: %+v", err)
	}
	o.Configure(tenantConfigurationsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DashboardsClient:           dashboardsClient,
		TenantConfigurationsClient: tenantConfigurationsClient,
	}, nil
}
