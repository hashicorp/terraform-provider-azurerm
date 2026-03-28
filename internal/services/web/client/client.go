// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppServiceCertificateOrdersClient *appservicecertificateorders.AppServiceCertificateOrdersClient
	CertificatesClient                *certificates.CertificatesClient
	WebAppsClient                     *webapps.WebAppsClient

	// azure-sdk-for-go
	AppServiceEnvironmentsClientV1 *web.AppServiceEnvironmentsClient
	AppServicePlansClientV1        *web.AppServicePlansClient
	AppServicesClientV1            *web.AppsClient
	BaseClientV1                   *web.BaseClient
	CertificatesClientV1           *web.CertificatesClient
	CertificatesOrderClientV1      *web.AppServiceCertificateOrdersClient
	StaticSitesClientV1            *web.StaticSitesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	appServiceCertificateOrdersClient, err := appservicecertificateorders.NewAppServiceCertificateOrdersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building App Service Certificate Orders client: %w", err)
	}
	o.Configure(appServiceCertificateOrdersClient.Client, o.Authorizers.ResourceManager)

	certificatesClient, err := certificates.NewCertificatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Certificates client: %w", err)
	}
	o.Configure(certificatesClient.Client, o.Authorizers.ResourceManager)

	webAppsClient, err := webapps.NewWebAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Web Apps client: %w", err)
	}
	o.Configure(webAppsClient.Client, o.Authorizers.ResourceManager)

	// Track 1
	appServiceEnvironmentsClient := web.NewAppServiceEnvironmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServiceEnvironmentsClient.Client, o.ResourceManagerAuthorizer)

	appServicePlansClient := web.NewAppServicePlansClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicePlansClient.Client, o.ResourceManagerAuthorizer)

	appServicesClient := web.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appServicesClient.Client, o.ResourceManagerAuthorizer)

	baseClient := web.NewWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&baseClient.Client, o.ResourceManagerAuthorizer)

	certificatesClientV1 := web.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClientV1.Client, o.ResourceManagerAuthorizer)

	certificatesOrderClient := web.NewAppServiceCertificateOrdersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesOrderClient.Client, o.ResourceManagerAuthorizer)

	staticSitesClient := web.NewStaticSitesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&staticSitesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppServiceCertificateOrdersClient: appServiceCertificateOrdersClient,
		CertificatesClient:                certificatesClient,
		WebAppsClient:                     webAppsClient,

		AppServiceEnvironmentsClientV1: &appServiceEnvironmentsClient,
		AppServicePlansClientV1:        &appServicePlansClient,
		AppServicesClientV1:            &appServicesClient,
		BaseClientV1:                   &baseClient,
		CertificatesClientV1:           &certificatesClientV1,
		CertificatesOrderClientV1:      &certificatesOrderClient,
		StaticSitesClientV1:            &staticSitesClient,
	}, nil
}
