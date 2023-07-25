// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FunctionsClient       *functions.FunctionsClient
	JobsClient            *streamingjobs.StreamingJobsClient
	InputsClient          *inputs.InputsClient
	OutputsClient         *outputs.OutputsClient
	TransformationsClient *transformations.TransformationsClient
	ClustersClient        *clusters.ClustersClient
	EndpointsClient       *privateendpoints.PrivateEndpointsClient
}

func NewClient(o *common.ClientOptions) *Client {
	functionsClient := functions.NewFunctionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&functionsClient.Client, o.ResourceManagerAuthorizer)

	jobsClient := streamingjobs.NewStreamingJobsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&jobsClient.Client, o.ResourceManagerAuthorizer)

	inputsClient := inputs.NewInputsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&inputsClient.Client, o.ResourceManagerAuthorizer)

	outputsClient := outputs.NewOutputsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&outputsClient.Client, o.ResourceManagerAuthorizer)

	transformationsClient := transformations.NewTransformationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&transformationsClient.Client, o.ResourceManagerAuthorizer)

	clustersClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clustersClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := privateendpoints.NewPrivateEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
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
