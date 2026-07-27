package netappaccounts

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetAppAccountsClient struct {
	Client *resourcemanager.Client
}

func NewNetAppAccountsClientWithBaseURI(sdkApi sdkEnv.Api) (*NetAppAccountsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "netappaccounts", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NetAppAccountsClient: %+v", err)
	}

	return &NetAppAccountsClient{
		Client: client,
	}, nil
}
