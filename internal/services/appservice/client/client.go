// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type Client struct {
	AppServiceEnvironmentClient *web.AppServiceEnvironmentsClient
	BaseClient                  *web.BaseClient
	ServicePlanClient           *web.AppServicePlansClient
	WebAppsClient               *web.AppsClient
}

func NewClient(o *common.ClientOptions) *Client {
	appServiceEnvironmentClient := web.NewAppServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	webAppServiceClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webAppServiceClient.Client, o.ResourceManagerAuthorizer)

	servicePlanClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicePlanClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServiceEnvironmentClient: &appServiceEnvironmentClient,
		BaseClient:                  &baseClient,
		ServicePlanClient:           &servicePlanClient,
		WebAppsClient:               &webAppServiceClient,
	}
}
