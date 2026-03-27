package rollingupgradestatusinfos

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RollingUpgradeStatusInfosClient struct {
	Client *resourcemanager.Client
}

func NewRollingUpgradeStatusInfosClientWithBaseURI(sdkApi sdkEnv.Api) (*RollingUpgradeStatusInfosClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "rollingupgradestatusinfos", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating RollingUpgradeStatusInfosClient: %+v", err)
	}

	return &RollingUpgradeStatusInfosClient{
		Client: client,
	}, nil
}
