package v2020_01_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverstart"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverstop"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2020-01-01/serverupgrade"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	ServerKeys    *serverkeys.ServerKeysClient
	ServerStart   *serverstart.ServerStartClient
	ServerStop    *serverstop.ServerStopClient
	ServerUpgrade *serverupgrade.ServerUpgradeClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	serverKeysClient, err := serverkeys.NewServerKeysClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerKeys client: %+v", err)
	}
	configureFunc(serverKeysClient.Client)

	serverStartClient, err := serverstart.NewServerStartClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerStart client: %+v", err)
	}
	configureFunc(serverStartClient.Client)

	serverStopClient, err := serverstop.NewServerStopClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerStop client: %+v", err)
	}
	configureFunc(serverStopClient.Client)

	serverUpgradeClient, err := serverupgrade.NewServerUpgradeClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ServerUpgrade client: %+v", err)
	}
	configureFunc(serverUpgradeClient.Client)

	return &Client{
		ServerKeys:    serverKeysClient,
		ServerStart:   serverStartClient,
		ServerStop:    serverStopClient,
		ServerUpgrade: serverUpgradeClient,
	}, nil
}
