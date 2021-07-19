package client

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-01-15/web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	WebAppsClient     *web.AppsClient
	ServicePlanClient *web.AppServicePlansClient
	BaseClient        *web.BaseClient
}

func NewClient(o *common.ClientOptions) *Client {
	appServiceClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceClient.Client, o.ResourceManagerAuthorizer)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	servicePlanClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicePlanClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BaseClient:        &baseClient,
		ServicePlanClient: &servicePlanClient,
		WebAppsClient:     &appServiceClient,
	}
}
