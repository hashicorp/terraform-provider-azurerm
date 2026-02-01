package volumegroups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeGroupsClient struct {
	Client *resourcemanager.Client
}

func NewVolumeGroupsClientWithBaseURI(sdkApi sdkEnv.Api) (*VolumeGroupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "volumegroups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating VolumeGroupsClient: %+v", err)
	}

	return &VolumeGroupsClient{
		Client: client,
	}, nil
}
