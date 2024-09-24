package virtualwans

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualWANsClient struct {
	Client *resourcemanager.Client
}

func NewVirtualWANsClientWithBaseURI(sdkApi sdkEnv.Api) (*VirtualWANsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "virtualwans", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VirtualWANsClient: %+v", err)
	}

	return &VirtualWANsClient{
		Client: client,
	}, nil
}
