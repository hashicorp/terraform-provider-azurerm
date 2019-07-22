package streamanalytics

import (
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	FunctionsClient       streamanalytics.FunctionsClient
	JobsClient            streamanalytics.StreamingJobsClient
	InputsClient          streamanalytics.InputsClient
	OutputsClient         streamanalytics.OutputsClient
	TransformationsClient streamanalytics.TransformationsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.FunctionsClient = streamanalytics.NewFunctionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FunctionsClient.Client, o.ResourceManagerAuthorizer)

	c.JobsClient = streamanalytics.NewStreamingJobsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.JobsClient.Client, o.ResourceManagerAuthorizer)

	c.InputsClient = streamanalytics.NewInputsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.InputsClient.Client, o.ResourceManagerAuthorizer)

	c.OutputsClient = streamanalytics.NewOutputsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.OutputsClient.Client, o.ResourceManagerAuthorizer)

	c.TransformationsClient = streamanalytics.NewTransformationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.TransformationsClient.Client, o.ResourceManagerAuthorizer)
	return &c
}
