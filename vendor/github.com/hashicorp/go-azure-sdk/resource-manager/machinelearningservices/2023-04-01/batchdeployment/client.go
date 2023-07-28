package batchdeployment

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchDeploymentClient struct {
	Client *resourcemanager.Client
}

func NewBatchDeploymentClientWithBaseURI(api environments.Api) (*BatchDeploymentClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "batchdeployment", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BatchDeploymentClient: %+v", err)
	}

	return &BatchDeploymentClient{
		Client: client,
	}, nil
}
