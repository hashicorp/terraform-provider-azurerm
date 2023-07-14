// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/dashboard"
	"github.com/hashicorp/go-azure-sdk/resource-manager/portal/2019-01-01-preview/tenantconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DashboardsClient           *dashboard.DashboardClient
	TenantConfigurationsClient *tenantconfiguration.TenantConfigurationClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	dashboardsClient, err := dashboard.NewDashboardClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dashboard client: %+v", err)
	}
	o.Configure(dashboardsClient.Client, o.Authorizers.ResourceManager)

	tenantConfigurationsClient, err := tenantconfiguration.NewTenantConfigurationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TenantConfiguration client: %+v", err)
	}
	o.Configure(tenantConfigurationsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DashboardsClient:           dashboardsClient,
		TenantConfigurationsClient: tenantConfigurationsClient,
	}, nil
}
