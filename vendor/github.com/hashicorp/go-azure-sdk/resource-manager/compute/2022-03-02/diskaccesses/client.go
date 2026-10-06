package diskaccesses

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskAccessesClient struct {
	Client *resourcemanager.Client
}

func NewDiskAccessesClientWithBaseURI(sdkApi sdkEnv.Api) (*DiskAccessesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "diskaccesses", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating DiskAccessesClient: %+v", err)
	}

	return &DiskAccessesClient{
		Client: client,
	}, nil
}
