package client

import (
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatasetClient                 *datafactory.DatasetsClient
	FactoriesClient               *datafactory.FactoriesClient
	IntegrationRuntimesClient     *datafactory.IntegrationRuntimesClient
	LinkedServiceClient           *datafactory.LinkedServicesClient
	ManagedPrivateEndpointsClient *datafactory.ManagedPrivateEndpointsClient
	ManagedVirtualNetworksClient  *datafactory.ManagedVirtualNetworksClient
	PipelinesClient               *datafactory.PipelinesClient
	TriggersClient                *datafactory.TriggersClient
}

func NewClient(o *common.ClientOptions) *Client {
	DatasetClient := datafactory.NewDatasetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatasetClient.Client, o.ResourceManagerAuthorizer)

	FactoriesClient := datafactory.NewFactoriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FactoriesClient.Client, o.ResourceManagerAuthorizer)

	IntegrationRuntimesClient := datafactory.NewIntegrationRuntimesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&IntegrationRuntimesClient.Client, o.ResourceManagerAuthorizer)

	LinkedServiceClient := datafactory.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LinkedServiceClient.Client, o.ResourceManagerAuthorizer)

	ManagedPrivateEndpointsClient := datafactory.NewManagedPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ManagedPrivateEndpointsClient.Client, o.ResourceManagerAuthorizer)

	ManagedVirtualNetworksClient := datafactory.NewManagedVirtualNetworksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ManagedVirtualNetworksClient.Client, o.ResourceManagerAuthorizer)

	PipelinesClient := datafactory.NewPipelinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PipelinesClient.Client, o.ResourceManagerAuthorizer)

	TriggersClient := datafactory.NewTriggersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TriggersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatasetClient:                 &DatasetClient,
		FactoriesClient:               &FactoriesClient,
		IntegrationRuntimesClient:     &IntegrationRuntimesClient,
		LinkedServiceClient:           &LinkedServiceClient,
		ManagedPrivateEndpointsClient: &ManagedPrivateEndpointsClient,
		ManagedVirtualNetworksClient:  &ManagedVirtualNetworksClient,
		PipelinesClient:               &PipelinesClient,
		TriggersClient:                &TriggersClient,
	}
}
