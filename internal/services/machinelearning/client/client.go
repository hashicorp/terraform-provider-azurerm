// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Datastore               *datastore.DatastoreClient
	MachineLearningComputes *machinelearningcomputes.MachineLearningComputesClient
	Workspaces              *workspaces.WorkspacesClient
	ManagedNetwork          *managednetwork.ManagedNetworkClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	datastoreClient, err := datastore.NewDatastoreClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Datastore client: %+v", err)
	}
	o.Configure(datastoreClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	computesClient, err := machinelearningcomputes.NewMachineLearningComputesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MachineLearningComputes client: %+v", err)
	}
	o.Configure(computesClient.Client, o.Authorizers.ResourceManager)

	managedNetworkClient, err := managednetwork.NewManagedNetworkClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building MachineLearningNetwork client: %+v", err)
	}
	o.Configure(managedNetworkClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MachineLearningComputes: computesClient,
		Datastore:               datastoreClient,
		Workspaces:              workspacesClient,
		ManagedNetwork:          managedNetworkClient,
	}, nil
}
