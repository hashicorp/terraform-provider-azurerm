package client

import (
	"github.com/Azure/azure-sdk-for-go/services/appplatform/mgmt/2020-07-01/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AppsClient               *appplatform.AppsClient
	CertificatesClient       *appplatform.CertificatesClient
	ConfigServersClient      *appplatform.ConfigServersClient
	CustomDomainsClient      *appplatform.CustomDomainsClient
	MonitoringSettingsClient *appplatform.MonitoringSettingsClient
	DeploymentsClient        *appplatform.DeploymentsClient
	ServicesClient           *appplatform.ServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	appsClient := appplatform.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appsClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := appplatform.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	configServersClient := appplatform.NewConfigServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configServersClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := appplatform.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	deploymentsClient := appplatform.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)

	monitoringSettingsClient := appplatform.NewMonitoringSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitoringSettingsClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := appplatform.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AppsClient:               &appsClient,
		CertificatesClient:       &certificatesClient,
		ConfigServersClient:      &configServersClient,
		CustomDomainsClient:      &customDomainsClient,
		DeploymentsClient:        &deploymentsClient,
		MonitoringSettingsClient: &monitoringSettingsClient,
		ServicesClient:           &servicesClient,
	}
}
