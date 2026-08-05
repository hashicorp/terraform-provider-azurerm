package backupautomaticandondemands

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupAutomaticAndOnDemandsClient struct {
	Client *resourcemanager.Client
}

func NewBackupAutomaticAndOnDemandsClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupAutomaticAndOnDemandsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backupautomaticandondemands", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupAutomaticAndOnDemandsClient: %+v", err)
	}

	return &BackupAutomaticAndOnDemandsClient{
		Client: client,
	}, nil
}
