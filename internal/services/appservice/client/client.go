// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appserviceplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

type Client struct {
	AppServiceEnvironmentClient *appserviceenvironments.AppServiceEnvironmentsClient
	BaseClient                  *web.BaseClient
	ResourceProvidersClient     *resourceproviders.ResourceProvidersClient
	ServicePlanClient           *appserviceplans.AppServicePlansClient
	WebAppsClient               *webapps.WebAppsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appServiceEnvironmentClient, err := appserviceenvironments.NewAppServiceEnvironmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AppServiceEnvironments client: %+v", err)
	}
	o.Configure(appServiceEnvironmentClient.Client, o.Authorizers.ResourceManager)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	webAppServiceClient, err := webapps.NewWebAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building WebApps client: %+v", err)
	}
	o.Configure(webAppServiceClient.Client, o.Authorizers.ResourceManager)

	resourceProvidersClient, err := resourceproviders.NewResourceProvidersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ResourceProviders client: %+v", err)
	}
	o.Configure(resourceProvidersClient.Client, o.Authorizers.ResourceManager)

	servicePlanClient, err := appserviceplans.NewAppServicePlansClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServicePlan client: %+v", err)
	}
	o.Configure(servicePlanClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AppServiceEnvironmentClient: appServiceEnvironmentClient,
		BaseClient:                  &baseClient,
		ResourceProvidersClient:     resourceProvidersClient,
		ServicePlanClient:           servicePlanClient,
		WebAppsClient:               webAppServiceClient,
	}, nil
}
