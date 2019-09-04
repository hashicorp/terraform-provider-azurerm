package web

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AppServicePlansClient *web.AppServicePlansClient
	AppServicesClient     *web.AppsClient
	CertificatesClient    *web.CertificatesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AppServicePlansClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AppServicePlansClient.Client, o.ResourceManagerAuthorizer)

	AppServicesClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AppServicesClient.Client, o.ResourceManagerAuthorizer)

	CertificatesClient := web.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CertificatesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServicePlansClient: &AppServicePlansClient,
		AppServicesClient:     &AppServicesClient,
		CertificatesClient:    &CertificatesClient,
	}
}
