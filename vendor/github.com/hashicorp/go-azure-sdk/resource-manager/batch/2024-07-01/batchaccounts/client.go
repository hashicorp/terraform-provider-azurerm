package batchaccounts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchAccountsClient struct {
	Client *resourcemanager.Client
}

func NewBatchAccountsClientWithBaseURI(sdkApi sdkEnv.Api) (*BatchAccountsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "batchaccounts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BatchAccountsClient: %+v", err)
	}

	return &BatchAccountsClient{
		Client: client,
	}, nil
}
