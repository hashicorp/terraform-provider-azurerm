package backupvaults

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupVaultsClient struct {
	Client *resourcemanager.Client
}

func NewBackupVaultsClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupVaultsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backupvaults", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupVaultsClient: %+v", err)
	}

	return &BackupVaultsClient{
		Client: client,
	}, nil
}
