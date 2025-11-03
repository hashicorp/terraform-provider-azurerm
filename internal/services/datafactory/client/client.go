// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedvirtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

type Client struct {
	Factories                 *factories.FactoriesClient
	Credentials               *credentials.CredentialsClient
	DataFlowClient            *dataflows.DataFlowsClient
	IntegrationRuntimesClient *integrationruntimes.IntegrationRuntimesClient
	ManagedPrivateEndpoints   *managedprivateendpoints.ManagedPrivateEndpointsClient
	ManagedVirtualNetworks    *managedvirtualnetworks.ManagedVirtualNetworksClient
	PipelinesClient           *pipelines.PipelinesClient

	// TODO: convert to using hashicorp/go-azure-sdk
	DatasetClient       *datafactory.DatasetsClient
	LinkedServiceClient *datafactory.LinkedServicesClient
	TriggersClient      *datafactory.TriggersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	credentialsClient, err := credentials.NewCredentialsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Credentials client: %+v", err)
	}
	o.Configure(credentialsClient.Client, o.Authorizers.ResourceManager)

	dataFlowClient, err := dataflows.NewDataFlowsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Flows client: %+v", err)
	}
	o.Configure(dataFlowClient.Client, o.Authorizers.ResourceManager)

	factoriesClient, err := factories.NewFactoriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Factories client: %+v", err)
	}
	o.Configure(factoriesClient.Client, o.Authorizers.ResourceManager)

	integrationRuntimesClient, err := integrationruntimes.NewIntegrationRuntimesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Integration Runtimes Client: %+v", err)
	}
	o.Configure(integrationRuntimesClient.Client, o.Authorizers.ResourceManager)

	managedPrivateEndpointsClient, err := managedprivateendpoints.NewManagedPrivateEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedPrivateEndpoints client: %+v", err)
	}
	o.Configure(managedPrivateEndpointsClient.Client, o.Authorizers.ResourceManager)

	managedVirtualNetworksClient, err := managedvirtualnetworks.NewManagedVirtualNetworksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ManagedVirtualNetworks client: %+v", err)
	}
	o.Configure(managedVirtualNetworksClient.Client, o.Authorizers.ResourceManager)

	// TODO: port the below operations to use `hashicorp/go-azure-sdk` in time
	DatasetClient := datafactory.NewDatasetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatasetClient.Client, o.ResourceManagerAuthorizer)

	LinkedServiceClient := datafactory.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LinkedServiceClient.Client, o.ResourceManagerAuthorizer)

	PipelinesClient, err := pipelines.NewPipelinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pipelines client: %+v", err)
	}
	o.Configure(PipelinesClient.Client, o.Authorizers.ResourceManager)

	TriggersClient := datafactory.NewTriggersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TriggersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Factories:                 factoriesClient,
		Credentials:               credentialsClient,
		DataFlowClient:            dataFlowClient,
		IntegrationRuntimesClient: integrationRuntimesClient,
		ManagedPrivateEndpoints:   managedPrivateEndpointsClient,
		ManagedVirtualNetworks:    managedVirtualNetworksClient,
		PipelinesClient:           PipelinesClient,

		// TODO: port to `hashicorp/go-azure-sdk`
		DatasetClient:       &DatasetClient,
		LinkedServiceClient: &LinkedServiceClient,
		TriggersClient:      &TriggersClient,
	}, nil
}
