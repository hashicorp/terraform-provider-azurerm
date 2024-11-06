package backupsv2

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupsV2Client struct {
	Client *resourcemanager.Client
}

func NewBackupsV2ClientWithBaseURI(sdkApi sdkEnv.Api) (*BackupsV2Client, error) {
	client, err := resourcemanager.NewClient(sdkApi, "backupsv2", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating BackupsV2Client: %+v", err)
	}

	return &BackupsV2Client{
		Client: client,
	}, nil
}
