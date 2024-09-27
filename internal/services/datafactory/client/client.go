// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/credentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedvirtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

type Client struct {
	Factories               *factories.FactoriesClient
	Credentials             *credentials.CredentialsClient
	ManagedPrivateEndpoints *managedprivateendpoints.ManagedPrivateEndpointsClient
	ManagedVirtualNetworks  *managedvirtualnetworks.ManagedVirtualNetworksClient

	// TODO: convert to using hashicorp/go-azure-sdk
	DataFlowClient            *datafactory.DataFlowsClient
	DatasetClient             *datafactory.DatasetsClient
	IntegrationRuntimesClient *datafactory.IntegrationRuntimesClient
	LinkedServiceClient       *datafactory.LinkedServicesClient
	PipelinesClient           *datafactory.PipelinesClient
	TriggersClient            *datafactory.TriggersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	factoriesClient, err := factories.NewFactoriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Factories client: %+v", err)
	}
	o.Configure(factoriesClient.Client, o.Authorizers.ResourceManager)

	credentialsClient, err := credentials.NewCredentialsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Factories client: %+v", err)
	}
	o.Configure(credentialsClient.Client, o.Authorizers.ResourceManager)

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
	dataFlowClient := datafactory.NewDataFlowsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dataFlowClient.Client, o.ResourceManagerAuthorizer)

	DatasetClient := datafactory.NewDatasetsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatasetClient.Client, o.ResourceManagerAuthorizer)

	IntegrationRuntimesClient := datafactory.NewIntegrationRuntimesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&IntegrationRuntimesClient.Client, o.ResourceManagerAuthorizer)

	LinkedServiceClient := datafactory.NewLinkedServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LinkedServiceClient.Client, o.ResourceManagerAuthorizer)

	PipelinesClient := datafactory.NewPipelinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PipelinesClient.Client, o.ResourceManagerAuthorizer)

	TriggersClient := datafactory.NewTriggersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TriggersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Factories:               factoriesClient,
		Credentials:             credentialsClient,
		ManagedPrivateEndpoints: managedPrivateEndpointsClient,
		ManagedVirtualNetworks:  managedVirtualNetworksClient,

		// TODO: port to `hashicorp/go-azure-sdk`
		DataFlowClient:            &dataFlowClient,
		DatasetClient:             &DatasetClient,
		IntegrationRuntimesClient: &IntegrationRuntimesClient,
		LinkedServiceClient:       &LinkedServiceClient,
		PipelinesClient:           &PipelinesClient,
		TriggersClient:            &TriggersClient,
	}, nil
}
