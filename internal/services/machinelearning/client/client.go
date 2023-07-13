// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ComputeClient    *machinelearningcomputes.MachineLearningComputesClient
	WorkspacesClient *workspaces.WorkspacesClient
	DatastoreClient  *datastore.DatastoreClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	computeClient, err := machinelearningcomputes.NewMachineLearningComputesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Compute Client: %+v", err)
	}
	o.Configure(computeClient.Client, o.Authorizers.ResourceManager)

	workspacesClient, err := workspaces.NewWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces Client: %+v", err)
	}
	o.Configure(workspacesClient.Client, o.Authorizers.ResourceManager)

	datastoreClient, err := datastore.NewDatastoreClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Datastore Client: %+v", err)
	}
	o.Configure(datastoreClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ComputeClient:    computeClient,
		WorkspacesClient: workspacesClient,
		DatastoreClient:  datastoreClient,
	}, nil
}
