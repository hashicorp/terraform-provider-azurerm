package cosmosdb

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CosmosDBClient struct {
	Client *resourcemanager.Client
}

func NewCosmosDBClientWithBaseURI(sdkApi sdkEnv.Api) (*CosmosDBClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "cosmosdb", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating CosmosDBClient: %+v", err)
	}

	return &CosmosDBClient{
		Client: client,
	}, nil
}
