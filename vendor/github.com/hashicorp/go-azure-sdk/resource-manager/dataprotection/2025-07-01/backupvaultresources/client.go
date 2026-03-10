package backupvaultresources

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupVaultResourcesClient struct {
	Client *resourcemanager.Client
}

func NewBackupVaultResourcesClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupVaultResourcesClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backupvaultresources", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupVaultResourcesClient: %+v", err)
	}

	return &BackupVaultResourcesClient{
		Client: client,
	}, nil
}
