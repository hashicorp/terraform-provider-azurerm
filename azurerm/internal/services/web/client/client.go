package client

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AppServiceEnvironmentsClient *web.AppServiceEnvironmentsClient
	AppServicePlansClient        *web.AppServicePlansClient
	AppServicesClient            *web.AppsClient
	BaseClient                   *web.BaseClient
	CertificatesClient           *web.CertificatesClient
	CertificatesOrderClient      *web.AppServiceCertificateOrdersClient
}

func NewClient(o *common.ClientOptions) *Client {
	appServiceEnvironmentsClient := web.NewAppServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceEnvironmentsClient.Client, o.ResourceManagerAuthorizer)

	appServicePlansClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicePlansClient.Client, o.ResourceManagerAuthorizer)

	appServicesClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicesClient.Client, o.ResourceManagerAuthorizer)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := web.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	certificatesOrderClient := web.NewAppServiceCertificateOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesOrderClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServiceEnvironmentsClient: &appServiceEnvironmentsClient,
		AppServicePlansClient:        &appServicePlansClient,
		AppServicesClient:            &appServicesClient,
		BaseClient:                   &baseClient,
		CertificatesClient:           &certificatesClient,
		CertificatesOrderClient:      &certificatesOrderClient,
	}
}
