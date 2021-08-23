package client

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
