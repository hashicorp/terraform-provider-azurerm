package client

import (
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2020-03-01/streamanalytics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FunctionsClient       *streamanalytics.FunctionsClient
	JobsClient            *streamanalytics.StreamingJobsClient
	InputsClient          *inputs.InputsClient
	OutputsClient         *outputs.OutputsClient
	TransformationsClient *streamanalytics.TransformationsClient
	ClustersClient        *streamanalytics.ClustersClient
	EndpointsClient       *streamanalytics.PrivateEndpointsClient
}

func NewClient(o *common.ClientOptions) *Client {
	functionsClient := streamanalytics.NewFunctionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&functionsClient.Client, o.ResourceManagerAuthorizer)

	jobsClient := streamanalytics.NewStreamingJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobsClient.Client, o.ResourceManagerAuthorizer)

	inputsClient := inputs.NewInputsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&inputsClient.Client, o.ResourceManagerAuthorizer)

	outputsClient := outputs.NewOutputsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&outputsClient.Client, o.ResourceManagerAuthorizer)

	transformationsClient := streamanalytics.NewTransformationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&transformationsClient.Client, o.ResourceManagerAuthorizer)

	clustersClient := streamanalytics.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&clustersClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := streamanalytics.NewPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FunctionsClient:       &functionsClient,
		JobsClient:            &jobsClient,
		InputsClient:          &inputsClient,
		OutputsClient:         &outputsClient,
		TransformationsClient: &transformationsClient,
		ClustersClient:        &clustersClient,
		EndpointsClient:       &endpointsClient,
	}
}
