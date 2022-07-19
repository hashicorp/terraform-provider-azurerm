package client

import (
	legacy "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-08-01/containerservice"
	legacyacr "github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2019-06-01-preview/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-03-02-preview/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2021-03-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AgentPoolsClient                  *containerservice.AgentPoolsClient
	ContainerRegistryAgentPoolsClient *containerregistry.AgentPoolsClient
	ContainerInstanceClient           *containerinstance.ContainerInstanceClient
	KubernetesClustersClient          *containerservice.ManagedClustersClient
	MaintenanceConfigurationsClient   *containerservice.MaintenanceConfigurationsClient
	RegistriesClient                  *containerregistry.RegistriesClient
	ReplicationsClient                *containerregistry.ReplicationsClient
	ServicesClient                    *legacy.ContainerServicesClient
	WebhooksClient                    *containerregistry.WebhooksClient
	TokensClient                      *containerregistry.TokensClient
	ScopeMapsClient                   *containerregistry.ScopeMapsClient
	TasksClient                       *legacyacr.TasksClient
	ConnectedRegistriesClient         *containerregistry.ConnectedRegistriesClient

	Environment azure.Environment
}

func NewClient(o *common.ClientOptions) *Client {
	registriesClient := containerregistry.NewRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&registriesClient.Client, o.ResourceManagerAuthorizer)

	registryAgentPoolsClient := containerregistry.NewAgentPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&registryAgentPoolsClient.Client, o.ResourceManagerAuthorizer)

	webhooksClient := containerregistry.NewWebhooksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webhooksClient.Client, o.ResourceManagerAuthorizer)

	replicationsClient := containerregistry.NewReplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&replicationsClient.Client, o.ResourceManagerAuthorizer)

	tokensClient := containerregistry.NewTokensClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tokensClient.Client, o.ResourceManagerAuthorizer)

	scopeMapsClient := containerregistry.NewScopeMapsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&scopeMapsClient.Client, o.ResourceManagerAuthorizer)

	tasksClient := legacyacr.NewTasksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tasksClient.Client, o.ResourceManagerAuthorizer)

	containerInstanceClient := containerinstance.NewContainerInstanceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&containerInstanceClient.Client, o.ResourceManagerAuthorizer)

	// AKS
	kubernetesClustersClient := containerservice.NewManagedClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&kubernetesClustersClient.Client, o.ResourceManagerAuthorizer)

	agentPoolsClient := containerservice.NewAgentPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentPoolsClient.Client, o.ResourceManagerAuthorizer)

	maintenanceConfigurationsClient := containerservice.NewMaintenanceConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&maintenanceConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := legacy.NewContainerServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	connectedRegistriesClient := containerregistry.NewConnectedRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&connectedRegistriesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AgentPoolsClient:                  &agentPoolsClient,
		ContainerRegistryAgentPoolsClient: &registryAgentPoolsClient,
		KubernetesClustersClient:          &kubernetesClustersClient,
		ContainerInstanceClient:           &containerInstanceClient,
		MaintenanceConfigurationsClient:   &maintenanceConfigurationsClient,
		RegistriesClient:                  &registriesClient,
		WebhooksClient:                    &webhooksClient,
		ReplicationsClient:                &replicationsClient,
		ServicesClient:                    &servicesClient,
		Environment:                       o.Environment,
		TokensClient:                      &tokensClient,
		ScopeMapsClient:                   &scopeMapsClient,
		TasksClient:                       &tasksClient,
		ConnectedRegistriesClient:         &connectedRegistriesClient,
	}
}
