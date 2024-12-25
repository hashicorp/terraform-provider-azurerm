package backups

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupsClient struct {
	Client *resourcemanager.Client
}

func NewBackupsClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupsClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backups", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupsClient: %+v", err)
	}

	return &BackupsClient{
		Client: client,
	}, nil
}
