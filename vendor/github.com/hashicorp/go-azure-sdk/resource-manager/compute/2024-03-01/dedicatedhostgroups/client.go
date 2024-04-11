package dedicatedhostgroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostGroupsClient struct {
	Client *resourcemanager.Client
}

func NewDedicatedHostGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*DedicatedHostGroupsClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "dedicatedhostgroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DedicatedHostGroupsClient: %+v", err)
	}

	return &DedicatedHostGroupsClient{
		Client: client,
	}, nil
}
