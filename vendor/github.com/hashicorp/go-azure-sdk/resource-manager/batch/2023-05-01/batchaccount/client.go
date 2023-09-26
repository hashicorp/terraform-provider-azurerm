package batchaccount

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchAccountClient struct {
	Client *resourcemanager.Client
}

func NewBatchAccountClientWithBaseURI(sdkApi sdkEnv.Api) (*BatchAccountClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "batchaccount", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BatchAccountClient: %+v", err)
	}

	return &BatchAccountClient{
		Client: client,
	}, nil
}
