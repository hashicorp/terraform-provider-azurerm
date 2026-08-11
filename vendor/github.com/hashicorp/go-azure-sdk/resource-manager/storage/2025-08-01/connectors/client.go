package connectors

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectorsClient struct {
	Client *resourcemanager.Client
}

func NewConnectorsClientWithBaseURI(sdkApi sdkEnv.Api) (*ConnectorsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "connectors", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ConnectorsClient: %+v", err)
	}

	return &ConnectorsClient{
		Client: client,
	}, nil
}
