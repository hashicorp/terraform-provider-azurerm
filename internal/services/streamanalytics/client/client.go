// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/functions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
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

func NewClient(o *common.ClientOptions) (*Client, error) {
	functionsClient, err := functions.NewFunctionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Functions client : %+v", err)
	}
	o.Configure(functionsClient.Client, o.Authorizers.ResourceManager)

	jobsClient, err := streamingjobs.NewStreamingJobsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Jobs client : %+v", err)
	}
	o.Configure(jobsClient.Client, o.Authorizers.ResourceManager)

	inputsClient, err := inputs.NewInputsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Inputs client : %+v", err)
	}
	o.Configure(inputsClient.Client, o.Authorizers.ResourceManager)

	outputsClient, err := outputs.NewOutputsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Outputs client : %+v", err)
	}
	o.Configure(outputsClient.Client, o.Authorizers.ResourceManager)

	transformationsClient, err := transformations.NewTransformationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Transformations client : %+v", err)
	}
	o.Configure(transformationsClient.Client, o.Authorizers.ResourceManager)

	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client : %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	endpointsClient, err := privateendpoints.NewPrivateEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Endpoints client : %+v", err)
	}
	o.Configure(endpointsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		FunctionsClient:       functionsClient,
		JobsClient:            jobsClient,
		InputsClient:          inputsClient,
		OutputsClient:         outputsClient,
		TransformationsClient: transformationsClient,
		ClustersClient:        clustersClient,
		EndpointsClient:       endpointsClient,
	}, nil
}
