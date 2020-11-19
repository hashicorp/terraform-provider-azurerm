package client

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2019-12-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	previewcontainerregistry "github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-06-01-preview/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2020-09-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AgentPoolsClient         *containerservice.AgentPoolsClient
	GroupsClient             *containerinstance.ContainerGroupsClient
	KubernetesClustersClient *containerservice.ManagedClustersClient
	RegistriesClient         *containerregistry.RegistriesClient
	TokensClient             *previewcontainerregistry.TokensClient
	ScopeMapsClient          *previewcontainerregistry.ScopeMapsClient
	ReplicationsClient       *containerregistry.ReplicationsClient
	ServicesClient           *containerservice.ContainerServicesClient
	WebhooksClient           *containerregistry.WebhooksClient

	Environment azure.Environment
}

func NewClient(o *common.ClientOptions) *Client {
	registriesClient := containerregistry.NewRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&registriesClient.Client, o.ResourceManagerAuthorizer)

	tokensClient := previewcontainerregistry.NewTokensClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tokensClient.Client, o.ResourceManagerAuthorizer)

	scopeMapsClient := previewcontainerregistry.NewScopeMapsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&scopeMapsClient.Client, o.ResourceManagerAuthorizer)

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

	servicesClient := containerservice.NewContainerServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AgentPoolsClient:         &agentPoolsClient,
		KubernetesClustersClient: &kubernetesClustersClient,
		GroupsClient:             &groupsClient,
		RegistriesClient:         &registriesClient,
		TokensClient:             &tokensClient,
		ScopeMapsClient:          &scopeMapsClient,
		WebhooksClient:           &webhooksClient,
		ReplicationsClient:       &replicationsClient,
		ServicesClient:           &servicesClient,
		Environment:              o.Environment,
	}
}
