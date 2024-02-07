package replicationvaultsetting

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationVaultSettingClient struct {
	Client *resourcemanager.Client
}

func NewReplicationVaultSettingClientWithBaseURI(sdkApi sdkEnv.Api) (*ReplicationVaultSettingClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "replicationvaultsetting", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating ReplicationVaultSettingClient: %+v", err)
	}

	return &ReplicationVaultSettingClient{
		Client: client,
	}, nil
}
