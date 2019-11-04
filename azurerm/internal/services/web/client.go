package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AppServicePlansClient        *web.AppServicePlansClient
	AppServicesClient            *web.AppsClient
	AppServiceEnvironmentsClient *web.AppServiceEnvironmentsClient
	CertificatesClient           *web.CertificatesClient
	CertificatesOrderClient      *web.AppServiceCertificateOrdersClient
	BaseClient                   *web.BaseClient
}

func BuildClient(o *common.ClientOptions) *Client {
	AppServicePlansClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AppServicePlansClient.Client, o.ResourceManagerAuthorizer)

	AppServiceEnvironmentsClient := web.NewAppServiceEnvironmentsClient(o.SubscriptionId)
	o.ConfigureClient(&AppServiceEnvironmentsClient.Client, o.ResourceManagerAuthorizer)

	AppServicesClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AppServicesClient.Client, o.ResourceManagerAuthorizer)

	CertificatesClient := web.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CertificatesClient.Client, o.ResourceManagerAuthorizer)

	CertificatesOrderClient := web.NewAppServiceCertificateOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CertificatesOrderClient.Client, o.ResourceManagerAuthorizer)

	BaseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BaseClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServicePlansClient:        &AppServicePlansClient,
		AppServiceEnvironmentsClient: &AppServiceEnvironmentsClient,
		AppServicesClient:            &AppServicesClient,
		CertificatesClient:           &CertificatesClient,
		CertificatesOrderClient:      &CertificatesOrderClient,
		BaseClient:                   &BaseClient,
	}
}
