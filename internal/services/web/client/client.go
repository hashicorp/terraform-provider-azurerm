// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AppServiceCertificateOrdersClient *appservicecertificateorders.AppServiceCertificateOrdersClient
	CertificatesClient                *certificates.CertificatesClient
	WebAppsClient                     *webapps.WebAppsClient
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

	return &Client{
		AppServiceCertificateOrdersClient: appServiceCertificateOrdersClient,
		CertificatesClient:                certificatesClient,
		WebAppsClient:                     webAppsClient,
	}, nil
}
