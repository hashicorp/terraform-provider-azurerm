package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2022-03-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	APIPortalClient             *appplatform.APIPortalsClient
	AppsClient                  *appplatform.AppsClient
	BindingsClient              *appplatform.BindingsClient
	BuildPackBindingClient      *appplatform.BuildpackBindingClient
	BuildServiceAgentPoolClient *appplatform.BuildServiceAgentPoolClient
	BuildServiceBuilderClient   *appplatform.BuildServiceBuilderClient
	CertificatesClient          *appplatform.CertificatesClient
	ConfigServersClient         *appplatform.ConfigServersClient
	ConfigurationServiceClient  *appplatform.ConfigurationServicesClient
	CustomDomainsClient         *appplatform.CustomDomainsClient
	GatewayClient               *appplatform.GatewaysClient
	GatewayCustomDomainClient   *appplatform.GatewayCustomDomainsClient
	GatewayRouteConfigClient    *appplatform.GatewayRouteConfigsClient
	MonitoringSettingsClient    *appplatform.MonitoringSettingsClient
	DeploymentsClient           *appplatform.DeploymentsClient
	ServicesClient              *appplatform.ServicesClient
	ServiceRegistryClient       *appplatform.ServiceRegistriesClient
	StoragesClient              *appplatform.StoragesClient
}

func NewClient(o *common.ClientOptions) *Client {
	apiPortalClient := appplatform.NewAPIPortalsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiPortalClient.Client, o.ResourceManagerAuthorizer)

	appsClient := appplatform.NewAppsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&appsClient.Client, o.ResourceManagerAuthorizer)

	bindingsClient := appplatform.NewBindingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&bindingsClient.Client, o.ResourceManagerAuthorizer)

	buildServiceAgentPoolClient := appplatform.NewBuildServiceAgentPoolClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&buildServiceAgentPoolClient.Client, o.ResourceManagerAuthorizer)

	buildpackBindingClient := appplatform.NewBuildpackBindingClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&buildpackBindingClient.Client, o.ResourceManagerAuthorizer)

	buildServiceBuilderClient := appplatform.NewBuildServiceBuilderClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&buildServiceBuilderClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := appplatform.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	configServersClient := appplatform.NewConfigServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configServersClient.Client, o.ResourceManagerAuthorizer)

	configurationServiceClient := appplatform.NewConfigurationServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationServiceClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := appplatform.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	deploymentsClient := appplatform.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)

	gatewayClient := appplatform.NewGatewaysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayClient.Client, o.ResourceManagerAuthorizer)

	gatewayCustomDomainClient := appplatform.NewGatewayCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayCustomDomainClient.Client, o.ResourceManagerAuthorizer)

	gatewayRouteConfigClient := appplatform.NewGatewayRouteConfigsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gatewayRouteConfigClient.Client, o.ResourceManagerAuthorizer)

	monitoringSettingsClient := appplatform.NewMonitoringSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitoringSettingsClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := appplatform.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	serviceRegistryClient := appplatform.NewServiceRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serviceRegistryClient.Client, o.ResourceManagerAuthorizer)

	storageClient := appplatform.NewStoragesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&storageClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		APIPortalClient:             &apiPortalClient,
		AppsClient:                  &appsClient,
		BindingsClient:              &bindingsClient,
		BuildPackBindingClient:      &buildpackBindingClient,
		BuildServiceAgentPoolClient: &buildServiceAgentPoolClient,
		BuildServiceBuilderClient:   &buildServiceBuilderClient,
		CertificatesClient:          &certificatesClient,
		ConfigServersClient:         &configServersClient,
		ConfigurationServiceClient:  &configurationServiceClient,
		CustomDomainsClient:         &customDomainsClient,
		DeploymentsClient:           &deploymentsClient,
		GatewayClient:               &gatewayClient,
		GatewayCustomDomainClient:   &gatewayCustomDomainClient,
		GatewayRouteConfigClient:    &gatewayRouteConfigClient,
		MonitoringSettingsClient:    &monitoringSettingsClient,
		ServicesClient:              &servicesClient,
		ServiceRegistryClient:       &serviceRegistryClient,
		StoragesClient:              &storageClient,
	}
}
