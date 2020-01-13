package client

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2018-09-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-10-01/containerservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AgentPoolsClient         *containerservice.AgentPoolsClient
	KubernetesClustersClient *containerservice.ManagedClustersClient
	GroupsClient             *containerinstance.ContainerGroupsClient
	RegistriesClient         *containerregistry.RegistriesClient
	WebhooksClient           *containerregistry.WebhooksClient
	ReplicationsClient       *containerregistry.ReplicationsClient

	// TODO: this can be removed in 2.0
	ServicesClient *containerservice.ContainerServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	registriesClient := containerregistry.NewRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&registriesClient.Client, o.ResourceManagerAuthorizer)

	webhooksClient := containerregistry.NewWebhooksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&webhooksClient.Client, o.ResourceManagerAuthorizer)

	replicationsClient := containerregistry.NewReplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&replicationsClient.Client, o.ResourceManagerAuthorizer)

	groupsClient := containerinstance.NewContainerGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&groupsClient.Client, o.ResourceManagerAuthorizer)

	// AKS
	kubernetesClustersClient := containerservice.NewManagedClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&kubernetesClustersClient.Client, o.ResourceManagerAuthorizer)

	agentPoolsClient := containerservice.NewAgentPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentPoolsClient.Client, o.ResourceManagerAuthorizer)

	// ACS
	servicesClient := containerservice.NewContainerServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AgentPoolsClient:         &agentPoolsClient,
		KubernetesClustersClient: &kubernetesClustersClient,
		GroupsClient:             &groupsClient,
		RegistriesClient:         &registriesClient,
		WebhooksClient:           &webhooksClient,
		ReplicationsClient:       &replicationsClient,
		ServicesClient:           &servicesClient,
	}
}
