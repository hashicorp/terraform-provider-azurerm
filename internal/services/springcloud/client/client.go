// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

type Client struct {
	APIPortalClient              *appplatform.APIPortalsClient
	APIPortalCustomDomainClient  *appplatform.APIPortalCustomDomainsClient
	ApplicationAcceleratorClient *appplatform.ApplicationAcceleratorsClient
	ApplicationLiveViewsClient   *appplatform.ApplicationLiveViewsClient
	AppsClient                   *appplatform.AppsClient
	BindingsClient               *appplatform.BindingsClient
	BuildPackBindingClient       *appplatform.BuildpackBindingClient
	BuildServiceAgentPoolClient  *appplatform.BuildServiceAgentPoolClient
	BuildServiceBuilderClient    *appplatform.BuildServiceBuilderClient
	BuildServiceClient           *appplatform.BuildServiceClient
	CertificatesClient           *appplatform.CertificatesClient
	ConfigServersClient          *appplatform.ConfigServersClient
	ConfigurationServiceClient   *appplatform.ConfigurationServicesClient
	ContainerRegistryClient      *appplatform.ContainerRegistriesClient
	CustomDomainsClient          *appplatform.CustomDomainsClient
	CustomizedAcceleratorClient  *appplatform.CustomizedAcceleratorsClient
	DevToolPortalClient          *appplatform.DevToolPortalsClient
	GatewayClient                *appplatform.GatewaysClient
	GatewayCustomDomainClient    *appplatform.GatewayCustomDomainsClient
	GatewayRouteConfigClient     *appplatform.GatewayRouteConfigsClient
	MonitoringSettingsClient     *appplatform.MonitoringSettingsClient
	DeploymentsClient            *appplatform.DeploymentsClient
	ServicesClient               *appplatform.ServicesClient
	ServiceRegistryClient        *appplatform.ServiceRegistriesClient
	StoragesClient               *appplatform.StoragesClient
}

func NewClient(o *common.ClientOptions) *Client {
	apiPortalClient := appplatform.NewAPIPortalsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiPortalClient.Client, o.ResourceManagerAuthorizer)

	apiPortalCustomDomainClient := appplatform.NewAPIPortalCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&apiPortalCustomDomainClient.Client, o.ResourceManagerAuthorizer)

	applicationAcceleratorClient := appplatform.NewApplicationAcceleratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationAcceleratorClient.Client, o.ResourceManagerAuthorizer)

	applicationLiveViewsClient := appplatform.NewApplicationLiveViewsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&applicationLiveViewsClient.Client, o.ResourceManagerAuthorizer)

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

	buildServiceClient := appplatform.NewBuildServiceClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&buildServiceClient.Client, o.ResourceManagerAuthorizer)

	certificatesClient := appplatform.NewCertificatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&certificatesClient.Client, o.ResourceManagerAuthorizer)

	configServersClient := appplatform.NewConfigServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configServersClient.Client, o.ResourceManagerAuthorizer)

	configurationServiceClient := appplatform.NewConfigurationServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationServiceClient.Client, o.ResourceManagerAuthorizer)

	containerRegistryClient := appplatform.NewContainerRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&containerRegistryClient.Client, o.ResourceManagerAuthorizer)

	customDomainsClient := appplatform.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	customizedAcceleratorClient := appplatform.NewCustomizedAcceleratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customizedAcceleratorClient.Client, o.ResourceManagerAuthorizer)

	deploymentsClient := appplatform.NewDeploymentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&deploymentsClient.Client, o.ResourceManagerAuthorizer)

	devToolPortalClient := appplatform.NewDevToolPortalsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&devToolPortalClient.Client, o.ResourceManagerAuthorizer)

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
		APIPortalClient:              &apiPortalClient,
		APIPortalCustomDomainClient:  &apiPortalCustomDomainClient,
		ApplicationAcceleratorClient: &applicationAcceleratorClient,
		ApplicationLiveViewsClient:   &applicationLiveViewsClient,
		AppsClient:                   &appsClient,
		BindingsClient:               &bindingsClient,
		BuildPackBindingClient:       &buildpackBindingClient,
		BuildServiceAgentPoolClient:  &buildServiceAgentPoolClient,
		BuildServiceBuilderClient:    &buildServiceBuilderClient,
		BuildServiceClient:           &buildServiceClient,
		CertificatesClient:           &certificatesClient,
		ConfigServersClient:          &configServersClient,
		ConfigurationServiceClient:   &configurationServiceClient,
		ContainerRegistryClient:      &containerRegistryClient,
		CustomDomainsClient:          &customDomainsClient,
		CustomizedAcceleratorClient:  &customizedAcceleratorClient,
		DeploymentsClient:            &deploymentsClient,
		DevToolPortalClient:          &devToolPortalClient,
		GatewayClient:                &gatewayClient,
		GatewayCustomDomainClient:    &gatewayCustomDomainClient,
		GatewayRouteConfigClient:     &gatewayRouteConfigClient,
		MonitoringSettingsClient:     &monitoringSettingsClient,
		ServicesClient:               &servicesClient,
		ServiceRegistryClient:        &serviceRegistryClient,
		StoragesClient:               &storageClient,
	}
}
