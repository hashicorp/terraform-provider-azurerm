package operationalizationclusters

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationalizationClustersClient struct {
	Client *resourcemanager.Client
}

func NewOperationalizationClustersClientWithBaseURI(api environments.Api) (*OperationalizationClustersClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(api, "operationalizationclusters", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating OperationalizationClustersClient: %+v", err)
	}

	return &OperationalizationClustersClient{
		Client: client,
	}, nil
}
