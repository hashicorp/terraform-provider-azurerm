package client

import (
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-02-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	WorkspacesClient             *machinelearningservices.WorkspacesClient
	MachineLearningComputeClient *machinelearningservices.MachineLearningComputeClient
	KubernetesClustersClient     *containerservice.ManagedClustersClient
	AgentPoolsClient             *containerservice.AgentPoolsClient
}

func NewClient(o *common.ClientOptions) *Client {
	WorkspacesClient := machinelearningservices.NewWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&WorkspacesClient.Client, o.ResourceManagerAuthorizer)

	MachineLearningComputeClient := machinelearningservices.NewMachineLearningComputeClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MachineLearningComputeClient.Client, o.ResourceManagerAuthorizer)

	KubernetesClustersClient := containerservice.NewManagedClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&KubernetesClustersClient.Client, o.ResourceManagerAuthorizer)

	agentPoolsClient := containerservice.NewAgentPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&agentPoolsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		WorkspacesClient:             &WorkspacesClient,
		MachineLearningComputeClient: &MachineLearningComputeClient,
		KubernetesClustersClient:     &KubernetesClustersClient,
		AgentPoolsClient:             &agentPoolsClient,
	}
}
