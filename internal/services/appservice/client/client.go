// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type Client struct {
	AppServiceEnvironmentClient *web.AppServiceEnvironmentsClient
	BaseClient                  *web.BaseClient
	ServicePlanClient           *web.AppServicePlansClient
	WebAppsClient               *web.AppsClient
	LinuxWebAppsClient          *webapps.WebAppsClient
	AvailabilityClient          *resourceproviders.ResourceProvidersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appServiceEnvironmentClient := web.NewAppServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceEnvironmentClient.Client, o.ResourceManagerAuthorizer)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	webAppServiceClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webAppServiceClient.Client, o.ResourceManagerAuthorizer)

	linuxWebAppServiceClient, err := webapps.NewWebAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WebApps client: %+v", err)
	}
	o.Configure(linuxWebAppServiceClient.Client, o.Authorizers.ResourceManager)

	servicePlanClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicePlanClient.Client, o.ResourceManagerAuthorizer)

	availabilityClient, err := resourceproviders.NewResourceProvidersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WebApps operation client: %+v", err)
	}
	o.Configure(availabilityClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AppServiceEnvironmentClient: &appServiceEnvironmentClient,
		BaseClient:                  &baseClient,
		ServicePlanClient:           &servicePlanClient,
		LinuxWebAppsClient:          linuxWebAppServiceClient,
		WebAppsClient:               &webAppServiceClient,
		AvailabilityClient:          availabilityClient,
	}, nil
}
