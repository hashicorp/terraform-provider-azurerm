package vipswap

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VipSwapClient struct {
	Client *resourcemanager.Client
}

func NewVipSwapClientWithBaseURI(sdkApi sdkEnv.Api) (*VipSwapClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "vipswap", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VipSwapClient: %+v", err)
	}

	return &VipSwapClient{
		Client: client,
	}, nil
}
