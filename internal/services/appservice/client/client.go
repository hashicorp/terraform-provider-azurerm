// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type Client struct {
	AppServiceEnvironmentClient *web.AppServiceEnvironmentsClient
	BaseClient                  *web.BaseClient
	ServicePlanClient           *web.AppServicePlansClient
	WebAppsClient               *webapps.WebAppsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appServiceEnvironmentClient := web.NewAppServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	webAppServiceClient, err := webapps.NewWebAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WebApps client: %+v", err)
	}
	o.Configure(webAppServiceClient.Client, o.Authorizers.ResourceManager)

	servicePlanClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicePlanClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServiceEnvironmentClient: &appServiceEnvironmentClient,
		BaseClient:                  &baseClient,
		ServicePlanClient:           &servicePlanClient,
		WebAppsClient:               webAppServiceClient,
	}, nil
}
