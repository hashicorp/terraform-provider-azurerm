// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppServiceEnvironmentsClientV1 *web.AppServiceEnvironmentsClient
	AppServicePlansClientV1        *web.AppServicePlansClient
	AppServicesClientV1            *web.AppsClient
	BaseClientV1                   *web.BaseClient
	CertificatesClientV1           *web.CertificatesClient
	CertificatesOrderClientV1      *web.AppServiceCertificateOrdersClient
	StaticSitesClientV1            *web.StaticSitesClient
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

	staticSitesClient := web.NewStaticSitesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&staticSitesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServiceEnvironmentsClientV1: &appServiceEnvironmentsClient,
		AppServicePlansClientV1:        &appServicePlansClient,
		AppServicesClientV1:            &appServicesClient,
		BaseClientV1:                   &baseClient,
		CertificatesClientV1:           &certificatesClient,
		CertificatesOrderClientV1:      &certificatesOrderClient,
		StaticSitesClientV1:            &staticSitesClient,
	}
}
