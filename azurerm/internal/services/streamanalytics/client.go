package streamanalytics

import (
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	FunctionsClient       *streamanalytics.FunctionsClient
	JobsClient            *streamanalytics.StreamingJobsClient
	InputsClient          *streamanalytics.InputsClient
	OutputsClient         *streamanalytics.OutputsClient
	TransformationsClient *streamanalytics.TransformationsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	FunctionsClient := streamanalytics.NewFunctionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FunctionsClient.Client, o.ResourceManagerAuthorizer)

	JobsClient := streamanalytics.NewStreamingJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&JobsClient.Client, o.ResourceManagerAuthorizer)

	InputsClient := streamanalytics.NewInputsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&InputsClient.Client, o.ResourceManagerAuthorizer)

	OutputsClient := streamanalytics.NewOutputsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&OutputsClient.Client, o.ResourceManagerAuthorizer)

	TransformationsClient := streamanalytics.NewTransformationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TransformationsClient.Client, o.ResourceManagerAuthorizer)
	return &Client{
		FunctionsClient:       &FunctionsClient,
		JobsClient:            &JobsClient,
		InputsClient:          &InputsClient,
		OutputsClient:         &OutputsClient,
		TransformationsClient: &TransformationsClient,
	}
}
