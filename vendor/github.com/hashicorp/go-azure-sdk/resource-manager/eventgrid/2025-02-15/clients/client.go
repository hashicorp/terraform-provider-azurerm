package clients

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClientsClient struct {
	Client *resourcemanager.Client
}

func NewClientsClientWithBaseURI(sdkApi sdkEnv.Api) (*ClientsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "clients", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ClientsClient: %+v", err)
	}

	return &ClientsClient{
		Client: client,
	}, nil
}
