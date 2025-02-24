package fileservice

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileServiceClient struct {
	Client *resourcemanager.Client
}

func NewFileServiceClientWithBaseURI(sdkApi sdkEnv.Api) (*FileServiceClient, error) {
	client, err := resourcemanager.NewClient(sdkApi, "fileservice", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating FileServiceClient: %+v", err)
	}

	return &FileServiceClient{
		Client: client,
	}, nil
}
