package serverupgrade

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerUpgradeClient struct {
	Client *resourcemanager.Client
}

func NewServerUpgradeClientWithBaseURI(sdkApi sdkEnv.Api) (*ServerUpgradeClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "serverupgrade", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ServerUpgradeClient: %+v", err)
	}

	return &ServerUpgradeClient{
		Client: client,
	}, nil
}
